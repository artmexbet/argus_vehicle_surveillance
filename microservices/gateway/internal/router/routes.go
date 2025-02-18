package router

import (
	"Argus/pkg/models"
	nats_connector "Argus/pkg/nats-connector"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

// AlarmOn ...
//
// @Summary Set car to security
// @Description Set car to security
// @ID alarm-on
// @Accept json
// @Produce json
// @Param body body models.AlarmOnRequest true "AlarmOnRequest"
// @Success 200 {object} models.AlarmOnResponse
// @Failure 400 {string} string "Cannot parse body"
// @Failure 400 {string} string "Struct is invalid"
// @Failure 500 {string} string "Cannot request message"
// @Router /alarm [post]
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

// GetCars ...
//
// @Summary Get cars
// @Description Returns list of cars which is in security
// @ID get-cars
// @Accept json
// @Produce json
// @Param login formData string true "Login"
// @Success 200 {array} models.SecurityCar
// @Failure 500 {string} string "Cannot marshal request"
// @Failure 500 {string} string "Cannot request message"
// @Router /cars [get]
func (r *Router) GetCars(ctx *fiber.Ctx) error {
	login := ctx.FormValue("login")
	req := &models.GetCarsRequest{Login: login}

	b, err := json.Marshal(req)
	if err != nil {
		slog.Error("Cannot marshal request", err.Error())
		return ctx.SendStatus(500)
	}

	request, err := r.broker.Request("get-cars", b)
	if err != nil {
		slog.Error("Cannot request message", err.Error())
		return ctx.SendStatus(500)
	}

	var resp nats_connector.NetworkResponse
	json.Unmarshal(request, &resp)
	ctx.SendStatus(resp.HTTPCode)
	return ctx.Send(resp.Data)
}
