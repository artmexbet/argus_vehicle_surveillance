package handler

import (
	"Argus/pkg/models"
	natsConnector "Argus/pkg/nats-connector"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"log/slog"
)

func (h *Handler) HandleAlarmOn() nats.MsgHandler {
	return func(msg *nats.Msg) {
		var alarmOnRequest models.AlarmOnRequest
		if err := json.Unmarshal(msg.Data, &alarmOnRequest); err != nil {
			slog.Error("Cant unmarshal data", err.Error())
			err := msg.Respond([]byte(err.Error()))
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}

		camId, _ := uuid.Parse("25d3e590-9870-11ef-a686-0242ac130002")
		accId, err := h.db.GetAccountIdByLogin(alarmOnRequest.Login)
		if err != nil {
			slog.Error(
				"Cannot get account with login",
				slog.String("login", alarmOnRequest.Login),
				slog.String("error", err.Error()),
			)

			bytes, _ := json.Marshal(natsConnector.NewResponse([]byte(err.Error()), 400))
			_ = msg.Respond(bytes)
			return
		}

		if has, err := h.db.CheckHasUserTelegramId(accId); err != nil {
			slog.Error(
				"Cannot check telegram id",
				slog.String("error", err.Error()),
				slog.Any("accountId", accId),
			)

			bytes, _ := json.Marshal(natsConnector.NewResponse([]byte(err.Error()), 500))
			_ = msg.Respond(bytes)
			return
		} else if !has {
			bytes, _ := json.Marshal(natsConnector.NewResponse([]byte("TelegramId is not specified"), 400))
			_ = msg.Respond(bytes)
			return
		}

		id, err := h.db.SetCarToSecurity(
			alarmOnRequest.CarID,
			models.CameraIDType(camId),
			accId,
			alarmOnRequest.Time,
		)
		if err != nil {
			slog.Error("Cannot set cars", err.Error())
			errResp := natsConnector.NewResponse([]byte(err.Error()), 500)
			bytes, _ := json.Marshal(errResp)
			err := msg.Respond(bytes)
			if err != nil {
				slog.Error(err.Error())
			}
			return
		}

		h.carProcessor.SetToSecurity(id)

		tmp := models.AlarmOnResponse{ID: id}
		data, err := json.Marshal(&tmp)
		if err != nil {
			err := msg.Respond([]byte(err.Error()))
			if err != nil {
				slog.Error(err.Error())
			}
			return
		}
		resp, _ := json.Marshal(natsConnector.NewResponse(data, 201))
		msg.Respond(resp)
	}
}

func (h *Handler) HandleGetCars() nats.MsgHandler {
	return func(msg *nats.Msg) {
		var req models.GetCarsRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			slog.Error("Cant unmarshal data", err.Error())
			err := msg.Respond([]byte(err.Error()))
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}

		cars, err := h.db.GetCarsByUserLogin(req.Login)
		if err != nil {
			slog.Error("Cannot get cars by user login", slog.String("login", req.Login), slog.String("error", err.Error()))
			bytes, _ := json.Marshal(natsConnector.NewResponse([]byte(err.Error()), 404))
			_ = msg.Respond(bytes)
			return
		}

		data, err := json.Marshal(cars)
		if err != nil {
			err := msg.Respond([]byte(err.Error()))
			if err != nil {
				slog.Error(err.Error())
			}
			return
		}
		resp, _ := json.Marshal(natsConnector.NewResponse(data, 200))
		err = msg.Respond(resp)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func (h *Handler) HandleAlarmOff() nats.MsgHandler {
	return func(msg *nats.Msg) {
		var alarmOffRequest models.AlarmOffRequest
		if err := json.Unmarshal(msg.Data, &alarmOffRequest); err != nil {
			slog.Error("Cant unmarshal data", err.Error())
			err := msg.Respond([]byte(err.Error()))
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}

		err := h.db.StopVehicleTracking(alarmOffRequest.ID)
		if err != nil {
			slog.Error("Cannot stop tracking", slog.Any("id", alarmOffRequest.ID), slog.String("error", err.Error()))
			bytes, _ := json.Marshal(natsConnector.NewResponse([]byte(err.Error()), 500))
			_ = msg.Respond(bytes)
			return
		}

		err = h.carProcessor.StopSecurity(alarmOffRequest.ID)
		if err != nil {
			slog.Error("Cannot stop security", slog.Any("id", alarmOffRequest.ID), slog.String("error", err.Error()))
			bytes, _ := json.Marshal(natsConnector.NewResponse([]byte(err.Error()), 500))
			_ = msg.Respond(bytes)
			return
		}

		resp, _ := json.Marshal(natsConnector.NewResponse([]byte("ok"), 200))
		err = msg.Respond(resp)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
