package service

import "Argus/pkg/models"

type Service struct {
	ch chan models.YOLOJson
}

func New() *Service {
	return &Service{
		ch: make(chan models.YOLOJson),
	}
}

func (s *Service) GetChannel() chan models.YOLOJson {
	return s.ch
}

func (s *Service) Send(msg models.YOLOJson) {
	s.ch <- msg
}
