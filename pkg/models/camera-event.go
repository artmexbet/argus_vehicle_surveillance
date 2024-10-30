package models

type CameraIDType int

type CameraEvent struct {
	ID        string       `json:"id"`
	CameraID  CameraIDType `json:"camera_id"`
	Timestamp string       `json:"timestamp"`
	EventType string       `json:"event_type"`
}
