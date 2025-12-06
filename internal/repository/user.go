package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kstsm/wb-event-booker/internal/apperrors"
	"github.com/kstsm/wb-event-booker/internal/models"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.conn.Exec(ctx, createUserQuery,
		user.ID,
		user.Name,
		user.Email,
		user.TelegramID,
		user.CreatedAt)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				switch pgError.ConstraintName {
				case "users_email_key":
					return apperrors.EmailAlreadyExists
				case "users_telegram_id_key":
					return apperrors.TelegramIDAlreadyExists
				}
			}
		}
		return fmt.Errorf("Exec-createUser: %w", err)
	}

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	var telegramID sql.NullInt64

	err := r.conn.QueryRow(ctx, getUserByIDQuery, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&telegramID,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.UserNotFound
		}
		return nil, fmt.Errorf("QueryRow-GetUserByID: %w", err)
	}

	if telegramID.Valid {
		user.TelegramID = &telegramID.Int64
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	var telegramID sql.NullInt64

	err := r.conn.QueryRow(ctx, getUserByEmailQuery, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&telegramID,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.UserNotFound
		}
		return nil, fmt.Errorf("QueryRow-GetUserByEmail: %w", err)

	}

	if telegramID.Valid {
		user.TelegramID = &telegramID.Int64
	}

	return &user, nil
}
