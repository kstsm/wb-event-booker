package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Date            time.Time `json:"date"`
	TotalSeats      int       `json:"total_seats"`
	ReservedSeats   int       `json:"reserved_seats"`
	BookedSeats     int       `json:"booked_seats"`
	BookingLifetime int       `json:"booking_lifetime"`
	PaymentReq      bool      `json:"requires_payment_confirmation"`
	CreatedAt       time.Time `json:"created_at"`
}

func (e *Event) AvailableSeats() int {
	return e.TotalSeats - e.ReservedSeats - e.BookedSeats
}
