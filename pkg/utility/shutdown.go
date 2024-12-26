package utility

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func EasyShutdown(
	waitSeconds int, component string, stopAction func() error,
) {
	EasyShutdownWithCtx(context.Background(), waitSeconds, component, stopAction)
}

func EasyShutdownWithCtx(
	countdown context.Context,
	waitSeconds int,
	component string,
	stopAction func() error,
) {
	NewShutdown(countdown, waitSeconds, slog.Default()).
		AddShutdownAction(component, stopAction).
		Serve()
}

// NewShutdown creates a new Shutdown instance that manages the graceful shutdown process.
//
// Parameters:
//
//   - countdown:
//     Specifies the context that determines when the graceful shutdown should be triggered.
//     If the context is canceled, the shutdown process will start.
//
//   - waitSeconds:
//     Specifies the maximum number of seconds to wait for the shutdown process to complete.
//     If this time elapses, the system will forcefully terminate regardless of the shutdown process's state.
//     A value <= 0 indicates it will wait permanently.
func NewShutdown(countdown context.Context, waitSeconds int, logger *slog.Logger) *Shutdown {
	osSig := make(chan os.Signal, 2)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancelCause(countdown)

	if logger == nil {
		logger = slog.Default()
	}

	shutdown := &Shutdown{
		osSig:     osSig,
		countdown: ctx,
		notify:    cancel,

		logger: logger,

		waitSeconds: waitSeconds,
		done:        make(chan struct{}),
	}

	return shutdown
}

type component struct {
	name string
	stop func() error
}

type Shutdown struct {
	osSig     chan os.Signal
	countdown context.Context
	notify    context.CancelCauseFunc

	waitSeconds int
	done        chan struct{}

	logger *slog.Logger
	mu     sync.Mutex

	// The fields `components` use an array of size 3,
	// representing three priority levels for shutdown process.
	//
	// priority 0 is the highest, and priority 2 is the lowest.
	componentQty int
	components   [3][]component
}

// AddPriorityShutdownAction registers a shutdown process with a given priority.
//
// Parameters:
//   - priority: Priority of the component (0 is the highest, and 2 is the lowest).
//   - name: Name of the component.
//   - stopAction: Function to execute during shutdown.
func (s *Shutdown) AddPriorityShutdownAction(priority uint, name string, stopAction func() error) *Shutdown {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		return s
	default:
	}

	s.componentQty++
	comp := component{name: name, stop: stopAction}
	s.components[priority] = append(s.components[priority], comp)
	return s
}

// AddShutdownAction This method registers a shutdown process to be stopped gracefully when a shutdown is triggered.
func (s *Shutdown) AddShutdownAction(name string, stopAction func() error) *Shutdown {
	const LowestPriority = len(s.components) - 1
	return s.AddPriorityShutdownAction(uint(LowestPriority), name, stopAction)
}

// Notify is used to trigger an immediate shutdown in case of a critical error.
func (s *Shutdown) Notify(cause error) {
	select {
	case <-s.done:
		return
	default:
		s.notify(cause)
	}
}

func (s *Shutdown) WaitChannel() <-chan struct{} {
	return s.done
}

func (s *Shutdown) Serve() {
	select {
	case sig := <-s.osSig:
		s.logger.Info("recv os signal",
			slog.String("trigger", "external"),
			slog.Any("signal", sig),
		)

	case <-s.countdown.Done():
		err := context.Cause(s.countdown)
		if errors.Is(err, context.Canceled) {
			s.logger.Info("recv go context",
				slog.String("trigger", "internal"),
			)
		} else {
			s.logger.Error("recv go context",
				slog.String("trigger", "internal"),
				slog.Any("err", err),
			)
		}

	case <-s.done:
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		return
	default:
	}

	defer close(s.done)

	s.logger.Info("shutdown start", slog.Int("qty", s.componentQty))
	start := time.Now()

	var timeout <-chan time.Time
	if s.waitSeconds > 0 {
		timeout = time.NewTimer(time.Duration(s.waitSeconds) * time.Second).C
	}

	finish := make(chan struct{}, 1)
	go func() {
		s.terminate()
		finish <- struct{}{}
	}()

	select {
	case <-timeout:
		duration := time.Since(start)
		s.logger.Error("shutdown failed because timeout", slog.String("duration", duration.String()))
	case <-finish:
		duration := time.Since(start)
		s.logger.Info("shutdown finish", slog.String("duration", duration.String()))
	}
}

func (s *Shutdown) terminate() {
	seq := 0
	wg := sync.WaitGroup{}
	for priority, components := range s.components {
		for j := range components {
			seq += 1
			wg.Add(1)
			comp := components[j]

			go func(sequence int) {
				defer wg.Done()

				logger := s.logger.With(
					slog.Int("no.", sequence),
					slog.Int("priority", priority),
					slog.String("component", comp.name),
				)

				logger.Info("terminate start")
				start := time.Now()

				err := comp.stop()
				if err != nil {
					logger.Error("terminate fail", slog.Any("err", err))
					return
				}

				duration := time.Since(start)
				logger.Info("terminate finish", slog.String("duration", duration.String()))

			}(seq)
		}
		wg.Wait()
	}
}
