package apperrors

import "errors"

var (
	EventNotFound              = errors.New("event not found")
	UserNotFound               = errors.New("user not found")
	NoAvailableSeats           = errors.New("no available seats")
	BookingNotFound            = errors.New("booking not found")
	BookingNotReserved         = errors.New("booking is not in reserved status")
	BookingDeadlinePassed      = errors.New("booking deadline has passed")
	UserAlreadyBookedThisEvent = errors.New("user already has a booking for this event")
	EventDoesNotRequirePayment = errors.New("event does not require payment confirmation")
	EventExpired               = errors.New("event has expired")
	EmailAlreadyExists         = errors.New("email already exists")
	TelegramIDAlreadyExists    = errors.New("telegram id already exists")
)
