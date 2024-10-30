package router

import "github.com/gofiber/fiber/v2"

func (r *Router) CameraList() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Camera list")
	}
}
