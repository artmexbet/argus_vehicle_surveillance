package router

import (
	"Argus/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) CameraList() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Camera list")
	}
}

func (r *Router) AlarmOn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(models.AlarmOnRequest)
		if err := ctx.BodyParser(req); err != nil {
			return ctx.SendStatus(400)
		}

		if err := r.validator.Struct(req); err != nil {
			return ctx.SendStatus(400)
		}

		return ctx.JSON(map[string]interface{}{"foo": "bar"})
	}
}
