package models

// SecurityCar is a struct that represents a security car
type SecurityCar struct {
	ID       SecurityCarIDType
	CameraID CameraIDType
	CarID    CarIDType
}

type AlarmOnRequest struct {
	CameraID  CameraIDType  `json:"camera_id" validate:"required"`
	AccountID AccountIDType `json:"account_id" validate:"required"`
	CarID     CarIDType     `json:"car_id" validate:"required"`
	Time      TimestampType `json:"time" validate:"required"`
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
