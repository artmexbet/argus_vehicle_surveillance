package service

type IBroker interface {
}

// Service ...
type Service struct {
}

// New ...
func New() *Service {
	return &Service{}
}

func (s *Service) RegisterCarAlarmSystem() {

}
