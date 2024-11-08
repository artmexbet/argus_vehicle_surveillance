package service

import (
	"Argus/pkg/models"
	"github.com/nats-io/nats.go"
)

func (s *Service) handleCamera(cameraID models.CameraIDType) nats.MsgHandler {
	return func(msg *nats.Msg) {
		err := s.ws.Publish("camera-01", msg.Data)
		if err != nil {
			return
		}

		// Do something with the message
	}
}
