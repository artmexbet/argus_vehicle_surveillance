package handler

import (
	"Argus/pkg/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"testing"
)

type Broker struct {
	Readers map[string]nats.MsgHandler
}

func NewBroker() *Broker {
	return &Broker{
		Readers: make(map[string]nats.MsgHandler),
	}
}

func (b *Broker) CreateReader(channel string, handler nats.MsgHandler) (*nats.Subscription, error) {
	b.Readers[channel] = handler
	return nil, nil
}

func (b *Broker) Publish(channel string, data []byte) error {
	if handler, ok := b.Readers[channel]; ok {
		msg := nats.NewMsg(channel)
		msg.Data = data
		handler(msg)
	}
	return nil
}

type WebSocket struct {
	Data map[string][][]byte
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		Data: make(map[string][][]byte),
	}
}

func (ws *WebSocket) Publish(channel string, data []byte) error {
	if _, ok := ws.Data[channel]; !ok {
		ws.Data[channel] = make([][]byte, 0)
	}
	ws.Data[channel] = append(ws.Data[channel], data)
	return nil
}

type Database struct {
	SecCars  map[models.CameraIDType][]models.SecurityCar
	Accounts map[string]models.AccountIDType
}

func NewDatabase(secCars map[models.CameraIDType][]models.SecurityCar, accounts map[string]models.AccountIDType) *Database {
	return &Database{
		SecCars:  secCars,
		Accounts: accounts,
	}
}

func (db *Database) GetAllSecuriedCarsByCamera(cameraID models.CameraIDType) ([]models.SecurityCar, error) {
	return db.SecCars[cameraID], nil
}

func (db *Database) CheckHasUserTelegramId(models.AccountIDType) (bool, error) {
	return true, nil
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

func (db *Database) GetAccountIdByLogin(login string) (models.AccountIDType, error) {
	return db.Accounts[login], nil
}

func TestHandler_HandleAlarmOn(t *testing.T) {
	broker := NewBroker()
	ws := NewWebSocket()
	db := NewDatabase(make(map[models.CameraIDType][]models.SecurityCar), make(map[string]models.AccountIDType))
	h := New(broker, ws, db, nil)
	err := h.Init()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	broker.Publish("alarm-on", []byte("test")) // TODO: test this
}

// TODO: Добавить моки для ошибок. Добавить тесты для ошибок.
