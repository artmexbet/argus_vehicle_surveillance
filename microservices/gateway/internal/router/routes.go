package router

import (
	"Argus/pkg/models"
	nats_connector "Argus/pkg/nats-connector"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func (r *Router) CameraList() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Camera list")
	}
}

func (r *Router) AlarmOn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(models.AlarmOnRequest)
		if err := json.Unmarshal(ctx.Body(), req); err != nil {
			slog.Error("Cannot parse body", err.Error())
			return ctx.SendStatus(400)
		}

		if err := r.validator.Struct(req); err != nil {
			slog.Error("Struct is invalid", err.Error(), *req)
			return ctx.SendStatus(400)
		}

		request, err := r.broker.Request("alarm-on", ctx.Body())
		if err != nil {
			slog.Error("Cannot request message", err.Error())
			return ctx.SendStatus(500)
		}
		var resp nats_connector.NetworkResponse
		json.Unmarshal(request, &resp)
		ctx.SendStatus(resp.HTTPCode)
		return ctx.Send(resp.Data)
	}
}
