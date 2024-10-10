package pkg

import (
	"net/http"
)

var (
	ErrInvalidParam = ErrorRegistry().
			AddErrorCode(4000).
			AddHttpStatus(http.StatusBadRequest).
			NewError("invalid parameter")
	ErrExists = ErrorRegistry().
			AddErrorCode(4001).
			AddHttpStatus(http.StatusConflict).
			NewError("resource already existed")
	ErrNotExists = ErrorRegistry().
			AddErrorCode(4002).
			AddHttpStatus(http.StatusNotFound).
			NewError("resource does not exist")

	ErrSystem = ErrorRegistry().
			AddErrorCode(5000).
			AddHttpStatus(http.StatusInternalServerError).
			NewError("system issue")
	ErrDatabase = ErrorRegistry().
			AddErrorCode(5001).
			AddHttpStatus(http.StatusInternalServerError).
			NewError("database issue")
)

var (
	ErrInvalidUsername = ErrorRegistry().
		AddErrorCode(6000).
		WrapError("username must be having a upper letter", ErrInvalidParam)
)
