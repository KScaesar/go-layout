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
			AddErrorCode(ErrCodeUndefined).
			Description("undefined error").
			HttpStatus(http.StatusInternalServerError).
			NewError()

	ErrInvalidParam = ErrorRegistry.
			AddErrorCode(4000).
			Description("invalid parameter").
			HttpStatus(http.StatusBadRequest).
			NewError()
	ErrExists = ErrorRegistry.
			AddErrorCode(4001).
			Description("resource already existed").
			HttpStatus(http.StatusConflict).
			NewError()
	ErrNotExists = ErrorRegistry.
			AddErrorCode(4002).
			Description("resource doesn't exist").
			HttpStatus(http.StatusNotFound).
			NewError()

	ErrSystem = ErrorRegistry.
			AddErrorCode(5000).
			Description("system issue").
			HttpStatus(http.StatusInternalServerError).
			NewError()
	ErrDatabase = ErrorRegistry.
			AddErrorCode(5001).
			Description("database issue").
			HttpStatus(http.StatusInternalServerError).
			NewError()
)
