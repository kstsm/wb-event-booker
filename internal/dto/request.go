package dto

type BookEventRequest struct {
	Email string `json:"email"`
}

type ConfirmBookingRequest struct {
	BookingID string `json:"booking_id"`
}

type CreateUserRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	TelegramID *int64 `json:"telegram_id,omitempty"`
}

type CreateEventRequest struct {
	Name                   string `json:"name"`
	Date                   string `json:"date"`
	TotalSeats             int    `json:"total_seats"`
	BookingLifetimeHours   int    `json:"booking_lifetime_hours"`
	BookingLifetimeMinutes int    `json:"booking_lifetime_minutes"`
	PaymentReq             bool   `json:"requires_payment_confirmation"`
}
