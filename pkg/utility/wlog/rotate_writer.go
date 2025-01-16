package wlog

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func openFileIfNotExist(FilePath string) (*os.File, error) {
	FilePath = filepath.Clean(FilePath)
	dir := filepath.Dir(FilePath)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create or open file: %w", err)
	}

	return file, nil
}

func NewRotateWriter(filename string, bufSize int) (io.WriteCloser, error) {
	file, err := openFileIfNotExist(filename)
	if err != nil {
		return nil, err
	}

	if bufSize <= 0 {
		const KB = 1 << 10
		bufSize = 16 * KB
	}

	conf := &Config{}
	conf.SetJsonFormat(true)
	logger := NewStderrLogger(conf)
	logger.WithAttribute(func(l *slog.Logger) *slog.Logger {
		return l.With(slog.String("file", filename))
	})

	w := &RotateWriter{
		filename: filename,
		bWriter:  bufio.NewWriterSize(file, bufSize),
		raw:      file,
		Logger:   logger.Slog(),
	}

	go w.autoFlush()
	go w.signalRotate()

	return w, nil
}

type RotateWriter struct {
	filename string
	bWriter  *bufio.Writer
	raw      *os.File
	isClosed bool
	mu       sync.Mutex
	Logger   *slog.Logger
}

func (w *RotateWriter) Write(p []byte) (nn int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.bWriter.Write(p)
}

func (w *RotateWriter) Close() (err error) {
	w.mu.Lock()
	defer func() {
		if err == nil {
			w.isClosed = true
		}
		w.mu.Unlock()
	}()

	if w.isClosed {
		return nil
	}

	if err = w.bWriter.Flush(); err != nil {
		w.Logger.Error("logger flush when close", "err", err)
	}
	return w.raw.Close()
}

// autoFlush
// To avoid high IOPS, reduce frequent disk write operations, and improve performance.
// while also preventing the situation where no new data is received, buffer is flushed at fixed intervals.
func (w *RotateWriter) autoFlush() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C

		w.mu.Lock()
		if w.isClosed {
			ticker.Stop()
			w.mu.Unlock()
			return
		}

		err := w.bWriter.Flush()
		if err != nil {
			w.Logger.Error("logger buffer auto flush", "err", err)
		}
		w.mu.Unlock()
	}
}

// signalRotate
// To integrate with external tools like `logrotate`, handling log rotation by listening for the `SIGHUP` signal.
// The application can switch log files without restarting, preventing log loss.
func (w *RotateWriter) signalRotate() {
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGHUP)
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-osSig:
			w.mu.Lock()
			if w.isClosed {
				ticker.Stop()
				w.mu.Unlock()
				return
			}

			file, err := openFileIfNotExist(w.filename)
			if err != nil {
				w.Logger.Error("create file when signal rotate", "err", err)
				w.mu.Unlock()
				break
			}

			if Err := w.bWriter.Flush(); Err != nil {
				w.Logger.Error("flush old file when signal rotate", "err", Err)
				file.Close()
				w.mu.Unlock()
				break
			}

			if Err := w.raw.Close(); Err != nil {
				w.Logger.Error("close old file When signal rotate", "err", Err)
				file.Close()
				w.mu.Unlock()
				break
			}

			w.Logger.Info("signal rotate success")
			*w.raw = *file
			w.mu.Unlock()

		case <-ticker.C:
			if w.isClosed {
				ticker.Stop()
				return
			}
		}
	}
}
