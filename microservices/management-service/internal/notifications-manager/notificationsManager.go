package notifications_manager

import (
	"Argus/pkg/models"
	"fmt"
	"log/slog"
)

// Config ...
type Config struct {
}

// NotificationsManager is a component which manages application notifications
type NotificationsManager struct {
	cfg *Config
}

// New creates a new notifications manager
func New(cfg *Config) *NotificationsManager {
	return &NotificationsManager{
		cfg: cfg,
	}
}

// SendNotification sends notification about a security car
func (nm *NotificationsManager) SendNotification(secId models.SecurityCarIDType, text string) error {
	slog.Info(fmt.Sprintf("Sending notification about security car %s: %s", secId, text))
	return nil
}
