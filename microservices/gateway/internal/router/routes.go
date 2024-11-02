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
		type Request struct {
			CameraID  models.CameraIDType  `json:"camera_id" validate:"required"`
			AccountID models.AccountIDType `json:"account_id" validate:"required"`
			CarID     models.CarIDType     `json:"car_id" validate:"required"`
			Time      string               `json:"time" validate:"required,datetime"`
		}

		req := new(Request)
		if err := ctx.BodyParser(req); err != nil {
			return ctx.SendStatus(400)
		}
		return ctx.JSON(map[string]interface{}{"foo": "bar"})
	}
}
