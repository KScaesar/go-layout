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

var DefaultShutdown *Shutdown

func SimpleShutdown(
	waitSeconds int, component string, stopAction func() error,
) {
	SimpleShutdownContext(context.Background(), waitSeconds, component, stopAction)
}

func SimpleShutdownContext(
	countdown context.Context,
	waitSeconds int,
	component string,
	stopAction func() error,
) {
	NewShutdown(countdown, waitSeconds).
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
//     A value of 0 means it will wait indefinitely.
func NewShutdown(countdown context.Context, waitSeconds int) *Shutdown {
	osSig := make(chan os.Signal, 2)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancelCause(countdown)

	shutdown := &Shutdown{
		osSig:     osSig,
		countdown: ctx,
		notify:    cancel,

		waitSeconds: waitSeconds,
		done:        make(chan struct{}),

		Logger: slog.Default(),
	}

	for i := range shutdown.names {
		shutdown.actions[i] = make(map[string]func() error)
	}

	return shutdown
}

type Shutdown struct {
	osSig     chan os.Signal
	countdown context.Context
	notify    context.CancelCauseFunc

	waitSeconds int
	done        chan struct{}

	Logger *slog.Logger
	mu     sync.Mutex

	// The fields `names`, `actions` and `waitBlocked` use an array of size 3,
	// representing three priority levels for shutdown process.
	//
	// priority 0 is the highest, and priority 2 is the lowest.
	actionsQty int
	names      [3][]string
	actions    [3]map[string]func() error
}

// AddPriorityShutdownAction registers a shutdown process with a given priority.
//
// Parameters:
//   - priority: Priority of the action (0 is the highest, and 2 is the lowest).
//   - component: Name of the components.
//   - stopAction: Function to execute during shutdown.
func (s *Shutdown) AddPriorityShutdownAction(priority uint, component string, stopAction func() error) *Shutdown {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		return s
	default:
	}

	s.actionsQty++
	s.names[priority] = append(s.names[priority], component)
	s.actions[priority][component] = stopAction
	return s
}

// AddShutdownAction This method registers a shutdown process to be stopped gracefully when a shutdown is triggered.
func (s *Shutdown) AddShutdownAction(component string, stopAction func() error) *Shutdown {
	const LowestPriority = len(s.names) - 1
	return s.AddPriorityShutdownAction(uint(LowestPriority), component, stopAction)
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
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		return
	default:
	}

	defer close(s.done)

	select {
	case sig := <-s.osSig:
		s.Logger.Info("recv os signal",
			slog.String("trigger", "external"),
			slog.Any("signal", sig),
		)

	case <-s.countdown.Done():
		err := context.Cause(s.countdown)
		if errors.Is(err, context.Canceled) {
			s.Logger.Info("recv go context",
				slog.String("trigger", "internal"),
			)
		} else {
			s.Logger.Error("recv go context",
				slog.String("trigger", "internal"),
				slog.Any("err", err),
			)
		}
	}

	s.Logger.Info("shutdown start", slog.Int("qty", s.actionsQty))
	start := time.Now()

	finish := make(chan struct{}, 1)
	var timeout <-chan time.Time
	if s.waitSeconds > 0 {
		timeout = time.NewTimer(time.Duration(s.waitSeconds) * time.Second).C
	}
	go func() {
		s.terminate()
		finish <- struct{}{}
	}()
	select {
	case <-timeout:
		s.Logger.Error("shutdown timeout")
	case <-finish:
	}

	duration := time.Since(start)
	s.Logger.Info("shutdown finish", slog.String("duration", duration.String()))
}

func (s *Shutdown) terminate() {
	seq := 0
	wg := sync.WaitGroup{}
	for i, components := range s.names {
		priority := i
		for j := range components {
			seq += 1
			wg.Add(1)
			component := components[j]
			go func(number int) {
				defer wg.Done()

				logger := s.Logger.With(
					slog.String("component", component),
					slog.Int("no.", number),
					slog.Int("priority", priority),
				)

				logger.Info("shutdown start")
				err := s.actions[priority][component]()
				if err != nil {
					logger.Error("shutdown fail", slog.Any("err", err))
					return
				}
				logger.Info("shutdown finish")

			}(seq)
		}
		wg.Wait()
	}
}
