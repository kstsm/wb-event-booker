package models

import (
	"github.com/google/uuid"
	"time"
)

type BookingStatus string

const (
	BookingStatusReserved  BookingStatus = "reserved"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID        uuid.UUID     `json:"id"`
	EventID   uuid.UUID     `json:"event_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Status    BookingStatus `json:"status"`
	Deadline  time.Time     `json:"deadline"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
