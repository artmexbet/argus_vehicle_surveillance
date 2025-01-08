package models

// SecurityCar is a struct that represents a security car
type SecurityCar struct {
	ID       SecurityCarIDType
	CameraID CameraIDType
	CarID    CarIDType
}

type AlarmOnRequest struct {
	Login string        `json:"login" validate:"required"`
	CarID CarIDType     `json:"car_id" validate:"required"`
	Time  TimestampType `json:"time" validate:"required"`
}

type AlarmOnResponse struct {
	ID SecurityCarIDType `json:"id"`
}

type YOLOObject struct {
	Id    int       `json:"id"`
	Class string    `json:"class"`
	BBox  []float32 `json:"bbox"`
}

type YOLOJson struct {
	FrameID int          `json:"frame_id"`
	Objects []YOLOObject `json:"objects"`
}
