package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"github.com/kstsm/wb-event-booker/internal/models"
	"github.com/kstsm/wb-event-booker/internal/repository"
)

type ServiceI interface {
	CreateEvent(ctx context.Context, req *dto.CreateEventRequest) (*models.Event, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	ListEvents(ctx context.Context) ([]*models.Event, error)
	BookEvent(ctx context.Context, eventID uuid.UUID, req *dto.BookEventRequest) (*dto.BookEventResponse, error)
	ConfirmBooking(ctx context.Context, eventID uuid.UUID, req *dto.ConfirmBookingRequest) error
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error)
	ListBookingsByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Booking, error)
}

type Service struct {
	repo repository.RepositoryI
}

func NewService(repo repository.RepositoryI) ServiceI {
	return &Service{
		repo: repo,
	}
}
