package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"github.com/kstsm/wb-event-booker/internal/models"
	"time"
)

func (s *Service) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		ID:         uuid.New(),
		Name:       req.Name,
		Email:      req.Email,
		TelegramID: req.TelegramID,
		CreatedAt:  time.Now(),
	}

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
