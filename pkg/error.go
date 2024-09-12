package pkg

import (
	"errors"
	"net/http"

	"github.com/KScaesar/go-layout/pkg/utility"
)

var ErrorRegistry = utility.NewErrorRegistry()

func ErrorUnwrap(err error) (myErr *utility.CustomError) {
	if errors.As(err, &myErr) {
		return myErr
	}
	return ErrUndefined.(*utility.CustomError)
}

const (
	ErrCodeUndefined = -1
)

var (
	ErrUndefined = ErrorRegistry.
			Register(ErrCodeUndefined).
			HttpStatus(http.StatusInternalServerError).
			Description("undefined error").
			Error()

	ErrInvalidParam = ErrorRegistry.
			Register(4000).
			HttpStatus(http.StatusBadRequest).
			Description("invalid parameter").
			Error()
	ErrExists = ErrorRegistry.
			Register(4001).
			HttpStatus(http.StatusConflict).
			Description("resource already existed").
			Error()
	ErrNotExists = ErrorRegistry.
			Register(4002).
			HttpStatus(http.StatusNotFound).
			Description("resource doesn't exist").
			Error()

	ErrSystem = ErrorRegistry.
			Register(5000).
			HttpStatus(http.StatusInternalServerError).
			Description("system issue").
			Error()
	ErrDatabase = ErrorRegistry.
			Register(5001).
			HttpStatus(http.StatusInternalServerError).
			Description("database issue").
			Error()
)
