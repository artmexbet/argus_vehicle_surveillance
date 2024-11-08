package service

import (
	"Argus/pkg/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var CameraID models.CameraIDType = models.CameraIDType(uuid.UUID{})

type Broker interface {
	CreateReader(string, nats.MsgHandler) (*nats.Subscription, error)
}

type WebSocket interface {
	Publish(string, []byte) error
}

type Service struct {
	broker Broker
	ws     WebSocket
}

func New(broker Broker, ws WebSocket) *Service {
	return &Service{
		broker: broker,
		ws:     ws,
	}
}

func (s *Service) Init() error {
	_, err := s.broker.CreateReader("camera-01", s.handleCamera(CameraID))
	if err != nil {
		return err
	}
	return nil
}
