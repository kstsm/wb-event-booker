package worker

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"github.com/kstsm/wb-event-booker/internal/models"
	"github.com/kstsm/wb-event-booker/internal/notifier"
	"github.com/kstsm/wb-event-booker/internal/repository"
)

type Worker struct {
	repo     repository.RepositoryI
	notifier notifier.NotifierI
}

func NewWorker(repo repository.RepositoryI, notifier notifier.NotifierI) *Worker {
	return &Worker{
		repo:     repo,
		notifier: notifier,
	}
}

func (w *Worker) ProcessExpiredBookings(ctx context.Context) error {
	slog.Info("Processing expired bookings...")

	bookings, err := w.repo.GetExpiredReservedBookings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired bookings: %w", err)
	}

	if len(bookings) == 0 {
		slog.Debug("No expired bookings found")
		return nil
	}

	slog.Infof("Found %d expired bookings to process", len(bookings))

	for _, booking := range bookings {
		if err := w.processExpiredBooking(ctx, booking); err != nil {
			slog.Error("Failed to process expired booking", "booking_id", booking.ID, "error", err)
			continue
		}
	}

	return nil
}

func (w *Worker) processExpiredBooking(ctx context.Context, booking *models.Booking) error {
	slog.Infof("Processing expired booking: booking_id=%s, event_id=%s, deadline=%v",
		booking.ID, booking.EventID, booking.Deadline)

	currentBooking, err := w.repo.GetBookingByID(ctx, booking.ID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if currentBooking.Status != models.BookingStatusReserved {
		slog.Infof("Booking %s is not in reserved status (status=%s), skipping", booking.ID, currentBooking.Status)
		return nil
	}

	err = w.cancelExpiredBooking(ctx, booking)
	if err != nil {
		return fmt.Errorf("failed to cancel expired booking: %w", err)
	}

	if w.notifier != nil {
		err = w.sendTelegramNotification(ctx, booking)
		if err != nil {
			slog.Warn("Failed to send Telegram notification", "booking_id", booking.ID, "error", err)
		}
	}

	slog.Infof("Successfully processed expired booking: booking_id=%s", booking.ID)
	return nil
}

func (w *Worker) cancelExpiredBooking(ctx context.Context, booking *models.Booking) error {
	err := w.repo.CancelExpiredBookingWithTransaction(ctx, booking.ID)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	slog.Infof("Cancelled expired booking: booking_id=%s, event_id=%s", booking.ID, booking.EventID)
	return nil
}

func (w *Worker) sendTelegramNotification(ctx context.Context, booking *models.Booking) error {
	user, err := w.repo.GetUserByID(ctx, booking.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	event, err := w.repo.GetEventByID(ctx, booking.EventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	if user.TelegramID == nil {
		slog.Infof("User %s has no Telegram ID, skipping Telegram notification", user.ID)
		return nil
	}

	message := fmt.Sprintf(
		dto.TelegramBookingCancel,
		event.Name,
		booking.ID,
		event.Name,
		event.Date.Format("2006-01-02 15:04"),
	)

	err = w.notifier.SendNotification(ctx, user.ID, *user.TelegramID, message)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
