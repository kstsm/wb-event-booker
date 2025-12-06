package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kstsm/wb-event-booker/internal/apperrors"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"github.com/kstsm/wb-event-booker/internal/models"
	"time"
)

func (s *Service) BookEvent(ctx context.Context,
	eventID uuid.UUID,
	req *dto.BookEventRequest,
) (*dto.BookEventResponse, error) {

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	booking, err := s.repo.BookEventWithTransaction(ctx, eventID, user.ID)
	if err != nil {
		return nil, err
	}

	resp := &dto.BookEventResponse{
		BookingID: booking.ID,
	}

	if !booking.Deadline.IsZero() {
		str := booking.Deadline.Format(time.RFC3339)
		resp.Deadline = &str
	}

	return resp, nil
}

func (s *Service) ConfirmBooking(ctx context.Context, eventID uuid.UUID, req *dto.ConfirmBookingRequest) error {
	booking, err := s.repo.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return err
	}

	if booking.EventID != eventID {
		return apperrors.BookingNotFound
	}

	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Date.Before(time.Now()) {
		return apperrors.EventExpired
	}

	if !event.PaymentReq {
		return apperrors.EventDoesNotRequirePayment
	}

	return s.repo.ConfirmBookingWithTransaction(ctx, req.BookingID)
}

func (s *Service) ListBookingsByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Booking, error) {
	return s.repo.GetBookingsByEventID(ctx, eventID)
}
