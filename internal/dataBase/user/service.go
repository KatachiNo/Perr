package user

import (
	"context"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type Service struct {
	storage Storage
	logger  *logg.Logger
}

func (s *Service) Create(ctx context.Context, dto CreateUserDTO) (User, error) {
	return User{}, nil
}
