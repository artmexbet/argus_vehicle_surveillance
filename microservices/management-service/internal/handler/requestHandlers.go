package handler

import (
	"Argus/pkg/models"
	"encoding/json"
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

		id, err := h.db.SetCarToSecurity(alarmOnRequest.CarID,
			alarmOnRequest.CameraID,
			alarmOnRequest.AccountID,
			alarmOnRequest.Time)
		if err != nil {
			slog.Error("Cannot get cars", err.Error())
			err := msg.Respond([]byte(err.Error()))
			if err != nil {
				slog.Error(err.Error())
			}
			return
		}

		type Response struct {
			ID models.SecurityCarIDType `json:"id"`
		}
		data, err := json.Marshal(Response{ID: id})
		if err != nil {
			err := msg.Respond([]byte(err.Error()))
			if err != nil {
				slog.Error(err.Error())
			}
			return
		}
		msg.Respond(data)
	}
}
