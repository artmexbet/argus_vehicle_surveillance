package handler

import (
	"Argus/pkg/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Broker struct {
	Readers map[string]nats.MsgHandler
}

func (b *Broker) CreateReader(channel string, handler nats.MsgHandler) (*nats.Subscription, error) {
	b.Readers[channel] = handler
	return nil, nil
}

type WebSocket struct {
	Data map[string][][]byte
}

func (ws *WebSocket) Publish(channel string, data []byte) error {
	if _, ok := ws.Data[channel]; !ok {
		ws.Data[channel] = make([][]byte, 0)
	}
	ws.Data[channel] = append(ws.Data[channel], data)
	return nil
}

type Database struct {
	SecCars map[models.CameraIDType][]models.SecurityCar
}

func (db *Database) GetAllSecuriedCarsByCamera(cameraID models.CameraIDType) ([]models.SecurityCar, error) {
	return db.SecCars[cameraID], nil
}

func (db *Database) SetCarToSecurity(
	carID models.CarIDType,
	cameraID models.CameraIDType,
	accountID models.AccountIDType,
	securityDateOn models.TimestampType,
) (models.SecurityCarIDType, error) {
	secCar := models.SecurityCar{
		ID:       models.SecurityCarIDType(uuid.New()),
		CameraID: cameraID,
		CarID:    carID,
	}
	db.SecCars[cameraID] = append(db.SecCars[cameraID], secCar)
	return secCar.ID, nil
}

// TODO: Добавить моки для ошибок. Добавить тесты для ошибок.
