package handler

import (
	"Argus/pkg/models"
	"github.com/nats-io/nats.go"
)

// handleCamera is a method that handles info about object from camera with cameraID
func (s *Service) handleCamera(cameraID models.CameraIDType) nats.MsgHandler {
	return func(msg *nats.Msg) {
		// Send message with secCar infos to centrifugo
		err := s.ws.Publish("camera-01", msg.Data)
		if err != nil {
			return
		}

		secCars, err := s.db.GetAllSecuriedCarsByCamera(cameraID)
		if err != nil {
			return
		}

		// Process cars. Sort it to security cars lists
		for _, secCar := range secCars {
			s.carProcessor.AppendCarInfo(secCar.ID, models.CarInfo{
				ID: secCar.CarID,
			})
		}
	}
}
