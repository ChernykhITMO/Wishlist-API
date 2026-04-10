package services

import (
	"context"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}
func (s *Service) Register(ctx context.Context, input RegisterInput) (string, error) {
	panic("implement")
}
