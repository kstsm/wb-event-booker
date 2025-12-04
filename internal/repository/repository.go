package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/wb-event-booker/internal/models"
)

type RepositoryI interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	ListEvents(ctx context.Context) ([]*models.Event, error)

	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	GetBookingsByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Booking, error)
	GetExpiredReservedBookings(ctx context.Context) ([]*models.Booking, error)

	CancelExpiredBookingWithTransaction(ctx context.Context, bookingID uuid.UUID) error
	ConfirmBookingWithTransaction(ctx context.Context, bookingID uuid.UUID) error
	BookEventWithTransaction(ctx context.Context, eventID, userID uuid.UUID) (*models.Booking, error)
}

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) RepositoryI {
	return &Repository{
		conn: conn,
	}
}
