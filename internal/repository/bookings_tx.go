package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-event-booker/internal/apperrors"
	"github.com/kstsm/wb-event-booker/internal/models"
	"time"
)

func (r *Repository) BookEventWithTransaction(ctx context.Context, eventID, userID uuid.UUID) (*models.Booking, error) {
	tx, err := r.conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("BeginTx-BookEventWithTransaction:: %w", err)
	}

	defer func() {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Errorf("Rollback-BookEventWithTransaction: %v", rbErr)
		}
	}()

	event, err := r.getEventForUpdate(ctx, tx, eventID)
	if err != nil {
		return nil, fmt.Errorf("getEventForUpdate-BookEventWithTransaction: %w", err)
	}

	if event.Date.Before(time.Now()) {
		return nil, apperrors.EventExpired
	}

	if event.TotalSeats-event.ReservedSeats-event.BookedSeats <= 0 {
		return nil, apperrors.NoAvailableSeats
	}

	exists, err := r.countUserBookings(ctx, tx, eventID, userID)
	if err != nil {
		return nil, fmt.Errorf("countUserBookings-BookEventWithTransaction: %w", err)
	}
	if exists > 0 {
		return nil, apperrors.UserAlreadyBookedThisEvent
	}

	status := models.BookingStatusConfirmed
	deadline := event.Date
	seatQuery := updateBookedSeatsQuery

	if event.PaymentReq {
		status = models.BookingStatusReserved
		deadline = time.Now().Add(time.Duration(event.BookingLifetime) * time.Minute).UTC()
		seatQuery = updateReservedSeatsQuery
	}

	booking := &models.Booking{
		ID:        uuid.New(),
		EventID:   eventID,
		UserID:    userID,
		Status:    status,
		Deadline:  deadline.UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if _, err = tx.Exec(ctx, seatQuery, eventID); err != nil {
		return nil, fmt.Errorf("updateSeats-BookEventWithTransaction: %w", err)
	}

	if err = r.insertBooking(ctx, tx, booking); err != nil {
		return nil, fmt.Errorf("insertBooking-BookEventWithTransaction: %w", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("Commit-BookEventWithTransaction:: %w", err)
	}

	return booking, nil
}

func (r *Repository) ConfirmBookingWithTransaction(ctx context.Context, bookingID uuid.UUID) error {
	tx, err := r.conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("BeginTx-ConfirmBookingWithTransaction: %w", err)
	}

	defer func() {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Errorf("Rollback-ConfirmBookingWithTransaction: %v", rbErr)
		}
	}()

	booking, err := r.getBookingInTx(ctx, tx, bookingID)
	if err != nil {
		return fmt.Errorf("getBookingInTx-ConfirmBookingWithTransaction: %w", err)
	}

	if booking.Status != models.BookingStatusReserved {
		return apperrors.BookingNotReserved
	}

	if time.Now().After(booking.Deadline) {
		return apperrors.BookingDeadlinePassed
	}

	if err = r.updateBookingStatus(ctx, tx, string(models.BookingStatusConfirmed), bookingID); err != nil {
		return fmt.Errorf("updateBookingStatus-ConfirmBookingWithTransaction: %w", err)
	}

	if err = r.updateEventSeatsReservedToBooked(ctx, tx, booking.EventID); err != nil {
		return fmt.Errorf("updateEventSeatsReservedToBooked-ConfirmBookingWithTransaction: %w", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("Commit-ConfirmBookingWithTransaction: %w", err)
	}

	return nil
}

func (r *Repository) CancelExpiredBookingWithTransaction(ctx context.Context, bookingID uuid.UUID) error {
	tx, err := r.conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("BeginTx-CancelExpiredBookingWithTransaction: %w", err)
	}

	defer func() {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Errorf("Rollback-CancelExpiredBookingWithTransaction: %v", rbErr)
		}
	}()

	booking, err := r.getBookingInTx(ctx, tx, bookingID)
	if err != nil {
		return fmt.Errorf("getBookingInTx-CancelExpiredBookingWithTransaction: %w", err)
	}

	if booking.Status != models.BookingStatusReserved {
		return apperrors.BookingNotReserved
	}

	if err = r.updateBookingStatus(ctx, tx, string(models.BookingStatusCancelled), bookingID); err != nil {
		return fmt.Errorf("updateBookingStatus-CancelExpiredBookingWithTransaction: %w", err)

	}

	if err = r.decreaseBookingSeats(ctx, tx, booking.EventID); err != nil {
		return err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("Commit-CancelExpiredBookingWithTransaction: %w", err)
	}

	return nil
}

func (r *Repository) getEventForUpdate(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) (*models.Event, error) {
	var event models.Event

	err := tx.QueryRow(ctx, selectEventForUpdateQuery, eventID).Scan(
		&event.ID,
		&event.Name,
		&event.Date,
		&event.TotalSeats,
		&event.ReservedSeats,
		&event.BookedSeats,
		&event.BookingLifetime,
		&event.PaymentReq,
		&event.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.EventNotFound
		}
		return nil, fmt.Errorf("QueryRow-getEventForUpdate: %w", err)
	}

	return &event, nil
}

func (r *Repository) updateEventSeatsReservedToBooked(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) error {
	_, err := tx.Exec(ctx, updateEventSeatsQuery, eventID)
	if err != nil {
		return fmt.Errorf("Exec-updateEventSeatsReservedToBooked: %w", err)
	}

	return nil
}

func (r *Repository) decreaseBookingSeats(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) error {
	_, err := tx.Exec(ctx, decreaseBookingSeatsQuery, eventID)
	if err != nil {
		return fmt.Errorf("Exec-decreaseBookingSeatsQuery: %w", err)
	}

	return nil
}

func (r *Repository) insertBooking(ctx context.Context, tx pgx.Tx, b *models.Booking) error {
	_, err := tx.Exec(ctx, insertBookingQuery,
		b.ID,
		b.EventID,
		b.UserID,
		b.Status,
		b.Deadline,
		b.CreatedAt,
		b.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("Exec-insertBooking: %w", err)
	}

	return nil
}

func (r *Repository) getBookingInTx(ctx context.Context, tx pgx.Tx, bookingID uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := tx.QueryRow(ctx, selectBookingForUpdateQuery, bookingID).Scan(
		&booking.ID,
		&booking.EventID,
		&booking.UserID,
		&booking.Status,
		&booking.Deadline,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.BookingNotFound
		}
		return nil, fmt.Errorf("QueryRow-getBookingInTx: %w", err)
	}

	return &booking, nil
}
