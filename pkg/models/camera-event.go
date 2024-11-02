package models

import "github.com/google/uuid"

type CameraIDType uuid.UUID
type CameraEventIDType uuid.UUID

type CameraEvent struct {
	ID        CameraEventIDType `json:"id"`
	CameraID  CameraIDType      `json:"camera_id"`
	Timestamp string            `json:"timestamp"`
	EventType string            `json:"event_type"`
}
