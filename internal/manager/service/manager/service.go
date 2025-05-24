package service

import "fmt"

// Publisher holds needed methods for the publisher client
type Publisher interface {
	PublishMessage(msg string) error
}

// Service holds the publisher client
type Service struct {
	publisher Publisher
}

// New returns a new service
func New(p Publisher) *Service {
	return &Service{
		publisher: p,
	}
}

// PublishMessage calls publisher in order to emit message
func (s *Service) PublishMessage(msg string) error {
	err := s.publisher.PublishMessage(msg)
	if err != nil {
		fmt.Println("there was an error publishing message")
	}
	return nil
}
