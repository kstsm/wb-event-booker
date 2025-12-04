package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"github.com/kstsm/wb-event-booker/internal/models"
	"time"
)

func (s *Service) CreateEvent(ctx context.Context, req *dto.CreateEventRequest) (*models.Event, error) {
	data, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	bookingLifetime := req.BookingLifetimeHours*60 + req.BookingLifetimeMinutes

	event := &models.Event{
		ID:              uuid.New(),
		Name:            req.Name,
		Date:            data.UTC(),
		TotalSeats:      req.TotalSeats,
		ReservedSeats:   0,
		BookedSeats:     0,
		BookingLifetime: bookingLifetime,
		PaymentReq:      req.PaymentReq,
		CreatedAt:       time.Now().UTC(),
	}

	err = s.repo.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Service) GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	return s.repo.GetEventByID(ctx, id)
}

func (s *Service) ListEvents(ctx context.Context) ([]*models.Event, error) {
	return s.repo.ListEvents(ctx)
}
