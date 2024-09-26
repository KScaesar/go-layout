package pkg

import (
	"net/http"

	"github.com/KScaesar/go-layout/pkg/utility"
)

func UnwrapError(err error) *utility.CustomError {
	myErr, ok := utility.UnwrapCustomError(err)
	if ok {
		return myErr
	}
	return ErrUndefined.(*utility.CustomError)
}

const (
	ErrCodeUndefined = -1
)

var (
	ErrUndefined = ErrorRegistry().
			AddErrorCode(ErrCodeUndefined).
			HttpStatus(http.StatusInternalServerError).
			NewError("undefined error")

	ErrInvalidParam = ErrorRegistry().
			AddErrorCode(4000).
			HttpStatus(http.StatusBadRequest).
			NewError("invalid parameter")
	ErrExists = ErrorRegistry().
			AddErrorCode(4001).
			HttpStatus(http.StatusConflict).
			NewError("resource already existed")
	ErrNotExists = ErrorRegistry().
			AddErrorCode(4002).
			HttpStatus(http.StatusNotFound).
			NewError("resource does not exist")

	ErrSystem = ErrorRegistry().
			AddErrorCode(5000).
			HttpStatus(http.StatusInternalServerError).
			NewError("system issue")
	ErrDatabase = ErrorRegistry().
			AddErrorCode(5001).
			HttpStatus(http.StatusInternalServerError).
			NewError("database issue")
)

var (
	ErrInvalidUsername = ErrorRegistry().
		AddErrorCode(6000).
		WrapError("username must be having a upper letter", ErrInvalidParam)
)
