package handler

import (
	"Argus/pkg/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var CameraID = models.CameraIDType(uuid.UUID{})

type IBroker interface {
	CreateReader(string, nats.MsgHandler) (*nats.Subscription, error)
}

type IWebSocket interface {
	Publish(string, []byte) error
}

type IDatabase interface {
	GetAllSecuriedCarsByCamera(models.CameraIDType) ([]models.SecurityCar, error)
	SetCarToSecurity(models.CarIDType, models.CameraIDType, models.AccountIDType) (models.SecurityCarIDType, error)
}

type ICarProcessor interface {
	AppendCarInfo(models.SecurityCarIDType, models.CarInfo)
	GetCarInfos(models.SecurityCarIDType) []models.CarInfo
}

type Handler struct {
	broker       IBroker
	ws           IWebSocket
	db           IDatabase
	carProcessor ICarProcessor
}

func New(broker IBroker, ws IWebSocket, db IDatabase, carProcessor ICarProcessor) *Handler {
	return &Handler{
		broker:       broker,
		ws:           ws,
		db:           db,
		carProcessor: carProcessor,
	}
}

func (h *Handler) Init() error {
	_, err := h.broker.CreateReader("camera-01", h.handleCamera(CameraID))
	if err != nil {
		return err
	}
	_, err = h.broker.CreateReader("alarm-on", h.handleAlarmOn())
	if err != nil {
		return err
	}
	return nil
}
