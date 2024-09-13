package pkg

import (
	"errors"
	"net/http"

	"github.com/KScaesar/go-layout/pkg/utility"
)

var defaultErrorRegistry = utility.NewErrorRegistry()

func ErrorRegistry() *utility.ErrorRegistry {
	return defaultErrorRegistry
}

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
	ErrUndefined = defaultErrorRegistry.
			AddErrorCode(ErrCodeUndefined).
			Description("undefined error").
			HttpStatus(http.StatusInternalServerError).
			NewError()

	ErrInvalidParam = defaultErrorRegistry.
			AddErrorCode(4000).
			Description("invalid parameter").
			HttpStatus(http.StatusBadRequest).
			NewError()
	ErrExists = defaultErrorRegistry.
			AddErrorCode(4001).
			Description("resource already existed").
			HttpStatus(http.StatusConflict).
			NewError()
	ErrNotExists = defaultErrorRegistry.
			AddErrorCode(4002).
			Description("resource doesn't exist").
			HttpStatus(http.StatusNotFound).
			NewError()

	ErrSystem = defaultErrorRegistry.
			AddErrorCode(5000).
			Description("system issue").
			HttpStatus(http.StatusInternalServerError).
			NewError()
	ErrDatabase = defaultErrorRegistry.
			AddErrorCode(5001).
			Description("database issue").
			HttpStatus(http.StatusInternalServerError).
			NewError()
)
