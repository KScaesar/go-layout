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

var DefaultShutdown = NewShutdown(context.Background(), 0)

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
	return &Shutdown{
		osSig:     osSig,
		countdown: ctx,
		notify:    cancel,

		waitSeconds: waitSeconds,
		done:        make(chan struct{}),

		Logger: slog.Default(),

		names:       make([]string, 0, 4),
		stopActions: make([]func() error, 0, 4),
	}
}

type Shutdown struct {
	osSig     chan os.Signal
	countdown context.Context
	notify    context.CancelCauseFunc

	waitSeconds int
	done        chan struct{}

	Logger *slog.Logger
	mu     sync.Mutex

	svcQty      int
	names       []string
	stopActions []func() error
}

// StopService This method registers a stop action to be stopped gracefully when a shutdown is triggered.
func (s *Shutdown) StopService(name string, stopAction func() error) *Shutdown {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		return s
	default:
	}

	fn := stopAction
	if s.waitSeconds > 0 {
		fn = func() error {
			result := make(chan error, 1)
			timeout := time.NewTimer(time.Duration(s.waitSeconds) * time.Second)

			go func() {
				result <- stopAction()
			}()

			select {
			case <-timeout.C:
				return errors.New("stop process timeout")
			case err := <-result:
				return err
			}
		}
	}

	s.svcQty++
	s.names = append(s.names, name)
	s.stopActions = append(s.stopActions, fn)
	return s
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

	// 有兩個觸發來源
	// 1. os.Signal 藉由系統外部觸發
	// 2. context countdown 藉由系統內部觸發
	select {
	case sig := <-s.osSig:
		s.Logger.Info("recv os signal: %v", sig)

	case <-s.countdown.Done():
		err := context.Cause(s.countdown)
		if errors.Is(err, context.Canceled) {
			s.Logger.Info("recv go context")
		} else {
			s.Logger.Error("recv go context: %v", err)
		}
	}

	s.Logger.Info("shutdown total service qty=%v", s.svcQty)
	start := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < s.svcQty; i++ {
		number := i + 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Logger.Info("number %v service %q shutdown start", number, s.names[number-1])
			err := s.stopActions[number-1]()
			if err != nil {
				s.Logger.Error("number %v service %q shutdown fail: %v", number, s.names[number-1], err)
				return
			}
			s.Logger.Info("number %v service %q shutdown finish", number, s.names[number-1])
		}()
	}
	wg.Wait()
	duration := time.Since(start)
	s.Logger.Info("shutdown finish", slog.String("duration", duration.String()))
}
