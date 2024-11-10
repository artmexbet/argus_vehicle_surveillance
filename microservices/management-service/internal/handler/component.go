package handler

import (
	"Argus/pkg/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var CameraID models.CameraIDType = models.CameraIDType(uuid.UUID{})

type IBroker interface {
	CreateReader(string, nats.MsgHandler) (*nats.Subscription, error)
}

type IWebSocket interface {
	Publish(string, []byte) error
}

type IDatabase interface {
	GetAllSecuriedCarsByCamera(models.CameraIDType) ([]models.SecurityCar, error)
}

type ICarProcessor interface {
	AppendCarInfo(models.SecurityCarIDType, models.CarInfo)
	GetCarInfos(models.SecurityCarIDType) []models.CarInfo
}

type Service struct {
	broker       IBroker
	ws           IWebSocket
	db           IDatabase
	carProcessor ICarProcessor
}

func New(broker IBroker, ws IWebSocket, db IDatabase, carProcessor ICarProcessor) *Service {
	return &Service{
		broker:       broker,
		ws:           ws,
		db:           db,
		carProcessor: carProcessor,
	}
}

func (s *Service) Init() error {
	_, err := s.broker.CreateReader("camera-01", s.handleCamera(CameraID))
	if err != nil {
		return err
	}
	return nil
}
