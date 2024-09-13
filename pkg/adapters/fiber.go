package adapters

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(c *fiber.Ctx, err error) error {
SendCustomError:
	myErr := pkg.ErrorUnwrap(err)
	if myErr.ErrorCode() != pkg.ErrCodeUndefined {
		return c.Status(myErr.HttpStatus()).JSON(fiber.Map{"msg": err.Error()})
	}

	isFixed, Err := fixUndefinedError(err)
	if isFixed {
		err = Err
		goto SendCustomError
	}

	pkg.Logger().CtxGetLogger(c.UserContext()).
		Warn("capture undefined error",
			slog.Any("err", err),
		)
	return c.Status(myErr.HttpStatus()).JSON(fiber.Map{"msg": err.Error()})
}

func fixUndefinedError(err error) (isFixed bool, Err error) {
	var fiberErr *fiber.Error

	switch {
	case errors.As(err, &fiberErr):
		switch fiberErr.Code {
		case fiber.StatusNotFound:
			return true, fmt.Errorf("%w: %w", pkg.ErrNotExists, err)
		}
	}

	return false, err
}
