package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-event-booker/internal/apperrors"
	"github.com/kstsm/wb-event-booker/internal/models"
)

func (r *Repository) CreateEvent(ctx context.Context, event *models.Event) error {
	_, err := r.conn.Exec(ctx, createEventQuery,
		event.ID,
		event.Name,
		event.Date,
		event.TotalSeats,
		event.BookingLifetime,
		event.PaymentReq,
		event.CreatedAt)
	if err != nil {
		return fmt.Errorf("Exec-CreateEvent: %w", err)
	}

	return nil
}

func (r *Repository) GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	var event models.Event

	err := r.conn.QueryRow(ctx, getEventByIDQuery, id).Scan(
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
		return nil, fmt.Errorf("QueryRow-GetEventByID: %w", err)
	}

	return &event, nil
}

func (r *Repository) ListEvents(ctx context.Context) ([]*models.Event, error) {
	rows, err := r.conn.Query(ctx, listEventsQuery)
	if err != nil {
		return nil, fmt.Errorf("Query-listEvents: %w", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event := new(models.Event)
		if err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Date,
			&event.TotalSeats,
			&event.ReservedSeats,
			&event.BookedSeats,
			&event.BookingLifetime,
			&event.PaymentReq,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan-listEvents: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Rows-listEvents: %w", err)
	}

	return events, nil
}
