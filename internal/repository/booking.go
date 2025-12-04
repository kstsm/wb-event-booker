package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-event-booker/internal/models"
)

func (r *Repository) GetBookingsByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Booking, error) {
	rows, err := r.conn.Query(ctx, getBookingsByEventQuery, eventID)
	if err != nil {
		return nil, fmt.Errorf("Query-GetBookingsByEventID query: %w", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		booking := new(models.Booking)
		if err := rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.UserID,
			&booking.Status,
			&booking.Deadline,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetBookingsByEventID scan: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBookingsByEventID rows.Err: %w", err)
	}

	return bookings, nil
}

func (r *Repository) GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	booking := new(models.Booking)
	err := r.conn.QueryRow(ctx, getBookingByIDQuery, id).Scan(
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
			return nil, fmt.Errorf("GetBookingByID: booking not found")
		}
		return nil, fmt.Errorf("QueryRow-GetBookingByID query: %w", err)
	}

	return booking, nil
}

func (r *Repository) GetExpiredReservedBookings(ctx context.Context) ([]*models.Booking, error) {
	rows, err := r.conn.Query(ctx, getExpiredReservedBookingsQuery)
	if err != nil {
		return nil, fmt.Errorf("Query-GetExpiredReservedBookings query: %w", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		booking := new(models.Booking)
		if err := rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.UserID,
			&booking.Status,
			&booking.Deadline,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetExpiredReservedBookings scan: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetExpiredReservedBookings rows.Err: %w", err)
	}

	return bookings, nil
}

func (r *Repository) updateBookingStatus(ctx context.Context, tx pgx.Tx, status string, bookingID uuid.UUID) error {
	_, err := tx.Exec(ctx, updateBookingStatusQuery, status, bookingID)
	if err != nil {
		return fmt.Errorf("Exec-updateBookingStatus: %w", err)
	}

	return nil
}

func (r *Repository) countUserBookings(ctx context.Context, tx pgx.Tx, eventID, userID uuid.UUID) (int, error) {
	var exists int
	err := tx.QueryRow(ctx, countUserBookingsQuery, eventID, userID).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("QueryRow-countUserBookings: %w", err)
	}

	return exists, nil
}
