package wlog

import (
	"bytes"
	"log/slog"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var DefaultFormats = []FormatFunc{
	FormatKeySource(),
	FormatKindTime,
	FormatKindDuration,
	FormatTypeFunc,
	FormatTypeStdError,
}

type FormatFunc func(groups []string, a slog.Attr) (slog.Attr, bool)

func FormatKeySource() FormatFunc {
	pool := sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}

	return func(groups []string, a slog.Attr) (slog.Attr, bool) {
		if a.Key == slog.SourceKey {
			src, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return a, false
			}

			buf := pool.Get().(*bytes.Buffer)
			defer func() {
				buf.Reset()
				pool.Put(buf)
			}()

			buf.WriteString(src.File)
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(src.Line))
			a.Value = slog.StringValue(buf.String())
			return a, true
		}
		return a, false
	}
}

func FormatKindTime(groups []string, a slog.Attr) (slog.Attr, bool) {
	if a.Value.Kind() == slog.KindTime {
		a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
		return a, true
	}
	return a, false
}

func FormatKindDuration(groups []string, a slog.Attr) (slog.Attr, bool) {
	if a.Value.Kind() == slog.KindDuration {
		a.Value = slog.StringValue(a.Value.Duration().String())
		return a, true
	}
	return a, false
}

func FormatTypeFunc(groups []string, a slog.Attr) (slog.Attr, bool) {
	if a.Value.Kind() == slog.KindAny {
		rv := reflect.ValueOf(a.Value.Any())
		if rv.Kind() != reflect.Func {
			return a, false
		}

		fnName := runtime.FuncForPC(rv.Pointer()).Name()
		a.Value = slog.StringValue(filepath.Base(fnName))
		return a, true
	}
	return a, false
}

func FormatTypeStdError(groups []string, a slog.Attr) (slog.Attr, bool) {
	if a.Value.Kind() == slog.KindAny {
		err, ok := a.Value.Any().(error)
		if !ok {
			return a, false
		}
		a.Value = slog.StringValue(err.Error())
		return a, true
	}
	return a, false
}
