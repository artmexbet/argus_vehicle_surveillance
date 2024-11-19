package models

import (
	"fmt"
	"github.com/google/uuid"
)

type CameraIDType uuid.UUID

func (cit *CameraIDType) UnmarshalJSON(b []byte) error {
	id, err := uuid.Parse(string(b[:]))
	if err != nil {
		return err
	}
	*cit = CameraIDType(id)
	return nil
}

func (cit *CameraIDType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", uuid.UUID(*cit).String())), nil
}

type CameraEventIDType uuid.UUID

type CameraEvent struct {
	ID        CameraEventIDType `json:"id"`
	CameraID  CameraIDType      `json:"camera_id"`
	Timestamp string            `json:"timestamp"`
	EventType string            `json:"event_type"`
}

type CarInfo struct {
	ID         CarIDType `json:"id"`
	Bbox       []float32 `json:"bbox"`
	IsCarFound bool      `json:"is_car_found"`
}
