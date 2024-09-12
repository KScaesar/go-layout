package dataflow

import (
	"fmt"
	"runtime/debug"
)

type ErrorHandleFunc func(message *Message, dep any, err error) error

type HandleFunc func(message *Message, dep any) error

func (h HandleFunc) PreMiddleware() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(message *Message, dep any) error {
			err := h(message, dep)
			if err != nil {
				return err
			}
			return next(message, dep)
		}
	}
}

func (h HandleFunc) PostMiddleware() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(message *Message, dep any) error {
			err := next(message, dep)
			if err != nil {
				return err
			}
			return h(message, dep)
		}
	}
}

func (h HandleFunc) Link(middlewares ...Middleware) HandleFunc {
	return Link(h, middlewares...)
}

type Middleware func(next HandleFunc) HandleFunc

func Link(handler HandleFunc, middlewares ...Middleware) HandleFunc {
	n := len(middlewares)
	for i := n - 1; 0 <= i; i-- {
		decorator := middlewares[i]
		handler = decorator(handler)
	}
	return handler
}

//

func UseRecover() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(message *Message, dep any) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("recovered from panic: %v\n%v", r, string(debug.Stack()))
				}
			}()
			return next(message, dep)
		}
	}
}
