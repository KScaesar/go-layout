package wlog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func NewBufioWriter(file *os.File, bufSize int, interval time.Duration) io.WriteCloser {
	w := &BufioWriter{
		Writer: bufio.NewWriterSize(file, bufSize),
		raw:    file,
	}
	if interval > 0 {
		go w.autoFlush(interval)
	}
	return w
}

type BufioWriter struct {
	*bufio.Writer
	raw      *os.File
	isClosed bool
	mu       sync.Mutex
}

func (w *BufioWriter) Fd() uintptr {
	return w.raw.Fd()
}

func (w *BufioWriter) Close() (err error) {
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

	if Err := w.Flush(); Err != nil {
		fmt.Printf("Error: BufioWriter Flush When Close: %v\n", Err)
	}
	return w.raw.Close()
}

func (w *BufioWriter) autoFlush(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C

		w.mu.Lock()
		if w.isClosed {
			ticker.Stop()
			w.mu.Unlock()
			return
		}

		err := w.Flush()
		if err != nil {
			fmt.Printf("Error: BufioWriter auto flush: %v\n", err)
		}
		w.mu.Unlock()
	}
}
