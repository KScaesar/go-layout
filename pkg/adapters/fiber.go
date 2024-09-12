package adapters

import (
	"github.com/KScaesar/go-layout/pkg"
	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	myErr := pkg.ErrorUnwrap(err)
	if myErr.CustomCode() == pkg.ErrCodeUndefined {
		pkg.Logger().CtxGetLogger(c.UserContext()).Warn("capture undefined error")
	}
	return c.Status(myErr.HttpStatus()).JSON(fiber.Map{"msg": err})
}
