package notifications_manager

import (
	"Argus/pkg/models"
	"encoding/json"
	"fmt"
	"log/slog"
)

// Config ...
type Config struct {
	NotificationTopic string `yaml:"topic" env-default:"security_car_notification"`
}

type IBroker interface {
	Publish(string, []byte) error
}

type IDatabase interface {
	GetTelegramId(models.SecurityCarIDType) (int64, error)
}

// NotificationsManager is a component which manages application notifications
type NotificationsManager struct {
	cfg    *Config
	broker IBroker
	db     IDatabase
}

// New creates a new notifications manager
func New(cfg *Config, broker IBroker, db IDatabase) *NotificationsManager {
	return &NotificationsManager{
		cfg:    cfg,
		broker: broker,
		db:     db,
	}
}

// SendNotification sends notification about a security car
func (nm *NotificationsManager) SendNotification(secId models.SecurityCarIDType, text string) error {
	slog.Info(fmt.Sprintf("Sending notification about security car %s: %s", secId, text))
	type SecurityCarNotification struct {
		UserId int64  `json:"telegram_id"`
		Text   string `json:"text"`
	}

	id, err := nm.db.GetTelegramId(secId)
	if err != nil {
		slog.Error("Cannot get telegram id", err.Error())
		return err
	}

	tmp := &SecurityCarNotification{UserId: id, Text: text}
	data, _ := json.Marshal(tmp)
	nm.broker.Publish(nm.cfg.NotificationTopic, data)
	return nil
}
