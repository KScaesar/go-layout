package adapters

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func HandleFiberError(c *fiber.Ctx, err error) error {
	myErr, ok := utility.UnwrapCustomError(err)
	if !ok {
		Err, isFixed := fixUnknownError(err)
		if isFixed {
			myErr, ok = utility.UnwrapCustomError(Err)
			err = Err
		}
	}

	// double check fix result
	if !ok {
		logger := pkg.Logger().CtxGetLogger(c.UserContext())
		logger.Warn("capture unknown error", slog.Any("err", err))
	}

	DefaultErrorResponse := fiber.Map{
		"code": myErr.ErrorCode(),
		"msg":  err.Error(),
	}
	return c.Status(myErr.HttpStatus()).JSON(DefaultErrorResponse)
}

func fixUnknownError(err error) (Err error, isFixed bool) {
	var fiberErr *fiber.Error

	switch {
	case errors.As(err, &fiberErr):
		Err, isFixed = fixFiberError(fiberErr)
		if isFixed {
			return Err, true
		}
	}

	return err, false
}

func fixFiberError(err *fiber.Error) (error, bool) {
	switch err.Code {
	case fiber.StatusNotFound:
		return fmt.Errorf("%w: %w", pkg.ErrNotExists, err), true
	}
	return err, false
}
