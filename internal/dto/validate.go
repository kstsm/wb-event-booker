package dto

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

const (
	MinBookingLifetime = 1
	MinTotalSeats      = 1
	MinTelegramID      = 1000000
	MaxTelegramID      = 9999999999
)

var (
	nameRegex  = regexp.MustCompile(`^[A-Za-zА-Яа-яЁё\s]+$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func (r *CreateUserRequest) ValidateUser() error {
	if r.Name == "" {
		return errors.New("user name is required")
	}

	if !nameRegex.MatchString(r.Name) {
		return errors.New("name must contain only letters")
	}

	if r.Email == "" {
		return errors.New("user email is required")
	}

	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}

	if r.TelegramID == nil {
		return nil
	}

	id := *r.TelegramID

	if id < MinTelegramID {
		return fmt.Errorf("telegram id must be >= %d", MinTelegramID)
	}

	if id > MaxTelegramID {
		return fmt.Errorf("telegram id must be <= %d", MaxTelegramID)
	}

	return nil
}

func (r *CreateEventRequest) ValidateEvent() error {
	if r.Name == "" {
		return errors.New("event name is required")
	}

	if r.TotalSeats < MinTotalSeats {
		return fmt.Errorf("total number of seats must be greater than or equal to %d", MinTotalSeats)
	}

	if r.BookingLifetimeHours < 0 {
		return errors.New("booking lifetime hours cannot be negative")
	}

	if r.BookingLifetimeMinutes < 0 || r.BookingLifetimeMinutes > 59 {
		return errors.New("booking lifetime minutes must be between 0 and 59")
	}

	bookingLifetime := r.BookingLifetimeHours*60 + r.BookingLifetimeMinutes
	if r.PaymentReq && bookingLifetime < MinBookingLifetime {
		return fmt.Errorf("minimum booking lifetime is %d minutes", MinBookingLifetime)
	}
	if !r.PaymentReq && bookingLifetime < 0 {
		return errors.New("booking lifetime cannot be negative")
	}

	data, err := time.Parse(time.RFC3339, r.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	if data.Before(time.Now().UTC()) {
		return errors.New("event date cannot be in the past")
	}

	return nil
}
