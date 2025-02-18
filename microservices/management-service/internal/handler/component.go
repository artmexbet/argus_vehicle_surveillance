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
	SetCarToSecurity(
		models.CarIDType,
		models.CameraIDType,
		models.AccountIDType,
		models.TimestampType,
	) (models.SecurityCarIDType, error)
	GetAccountIdByLogin(string) (models.AccountIDType, error)
	CheckHasUserTelegramId(models.AccountIDType) (bool, error)
	GetCarsByUserLogin(string) ([]models.SecurityCar, error)
}

type ICarProcessor interface {
	AppendCarInfo(models.SecurityCarIDType, models.CarInfo)
	GetCarInfos(models.SecurityCarIDType) []models.CarInfo
	SetToSecurity(models.SecurityCarIDType)
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
	tmp, _ := uuid.Parse("25d3e590-9870-11ef-a686-0242ac130002")
	CameraID = models.CameraIDType(tmp)
	_, err := h.broker.CreateReader("camera", h.HandleCamera(CameraID))
	if err != nil {
		return err
	}
	_, err = h.broker.CreateReader("alarm-on", h.HandleAlarmOn())
	if err != nil {
		return err
	}

	_, err = h.broker.CreateReader("get-cars", h.HandleGetCars())
	return nil
}
