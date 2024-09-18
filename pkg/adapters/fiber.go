package adapters

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	myErr := pkg.ErrorUnwrap(err)
	if myErr.ErrorCode() == pkg.ErrCodeUndefined {
		Err, isFixed := fixUndefinedError(err)
		if isFixed {
			err = Err
			myErr = pkg.ErrorUnwrap(err)
		} else {
			logger := pkg.Logger().CtxGetLogger(c.UserContext())
			logger.Warn("capture undefined error", slog.Any("err", err))
		}
	}
	return c.Status(myErr.HttpStatus()).JSON(fiber.Map{"msg": err.Error()})
}

func fixUndefinedError(err error) (Err error, isFixed bool) {
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
