package notifications_manager

import (
	"Argus/pkg/models"
	"fmt"
)

type Broker struct {
	Messages map[string][][]byte
}

func NewBroker() *Broker {
	return &Broker{Messages: make(map[string][][]byte)}
}

func (b *Broker) Publish(topic string, data []byte) error {
	if _, ok := b.Messages[topic]; !ok {
		b.Messages[topic] = make([][]byte, 0, 1)
	}
	b.Messages[topic] = append(b.Messages[topic], data)
	return nil
}

type Database struct {
	TelegramIds map[models.SecurityCarIDType]int64
}

func NewDatabase(telegramIds map[models.SecurityCarIDType]int64) *Database {
	return &Database{TelegramIds: telegramIds}
}

func (db *Database) GetTelegramId(secId models.SecurityCarIDType) (int64, error) {
	if _, ok := db.TelegramIds[secId]; !ok {
		return 0, fmt.Errorf("id attached to secCar with id %v not found", secId)
	}
	return db.TelegramIds[secId], nil
}
