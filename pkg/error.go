package pkg

import (
	"net/http"

	"github.com/KScaesar/go-layout/pkg/utility"
)

var ErrorRegistry = utility.NewErrorRegistry()

var (
	ErrUndefined = ErrorRegistry.
			Register(-1).
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
