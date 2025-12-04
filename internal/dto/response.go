package dto

import (
	"github.com/google/uuid"
	"github.com/kstsm/wb-event-booker/internal/models"
)

type CreateEventResponse struct {
	Event   *models.Event `json:"event"`
	Message string        `json:"message"`
}

type BookEventResponse struct {
	BookingID uuid.UUID `json:"booking_id"`
	Deadline  *string   `json:"deadline,omitempty"`
	Message   string    `json:"message"`
}

type ConfirmBookingResponse struct {
	Message string `json:"message"`
}

type GetEventResponse struct {
	Event *models.Event `json:"event"`
}

type ListEventsResponse struct {
	Events []*models.Event `json:"events"`
}

type ListBookingsResponse struct {
	Bookings []*models.Booking `json:"bookings"`
}

type CreateUserResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}
