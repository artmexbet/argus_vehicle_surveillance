package handler

import (
	"Argus/pkg/models"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log/slog"
)

// HandleCamera is a method that handles info about object from camera with cameraID
func (h *Handler) HandleCamera(cameraID models.CameraIDType) nats.MsgHandler {
	return func(msg *nats.Msg) {
		//slog.Info("Received message from camera %v", cameraID)
		// Send message with secCar infos to centrifugo
		err := h.ws.Publish("camera-01-ws", msg.Data)
		if err != nil {
			slog.Error("Error while sending message to centrifugo: %v", err)
			return
		}

		var yoloJson models.YOLOJson
		if err := json.Unmarshal(msg.Data, &yoloJson); err != nil {
			slog.Error("Error while unmarshalling message: %v", err)
			return
		}

		secCars, err := h.db.GetAllSecuriedCarsByCamera(cameraID)
		if err != nil {
			//slog.Error("Error while getting all securied cars by camera: %v", err)
			return
		}

		// Process cars. Sort it to security cars lists
		for _, secCar := range secCars {
			isCarFound := false
			for _, obj := range yoloJson.Objects {
				if obj.Class == "car" || obj.Class == "truck" {
					if secCar.CarID == models.CarIDType(obj.Id) {
						h.carProcessor.AppendCarInfo(secCar.ID, models.CarInfo{
							ID:         secCar.CarID,
							Bbox:       obj.BBox,
							IsCarFound: true,
						})
						isCarFound = true
						break
					}
				}
			}

			if !isCarFound {
				h.carProcessor.AppendCarInfo(secCar.ID, models.CarInfo{
					IsCarFound: false,
				})
			}
		}
	}
}
