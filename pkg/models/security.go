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
	Time      string        `json:"time" validate:"required,datetime"`
}
