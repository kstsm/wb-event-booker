package scheduler

import (
	"context"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-event-booker/internal/config"
	"time"
)

type SchedulerI interface {
	Start(ctx context.Context, cfg config.Config, task func() error)
	Stop()
}

type Scheduler struct {
	ticker *time.Ticker
	done   chan bool
}

func NewScheduler() SchedulerI {
	return &Scheduler{
		done: make(chan bool),
	}
}

func (s *Scheduler) Start(ctx context.Context, cfg config.Config, task func() error) {
	interval := time.Duration(cfg.Scheduler.CheckInterval) * time.Second

	s.ticker = time.NewTicker(interval)

	go func() {
		if err := task(); err != nil {
			slog.Error("Scheduler task error", "error", err)
		}

		for {
			select {
			case <-ctx.Done():
				slog.Info("Scheduler stopping due to context cancellation")
				s.Stop()
				return
			case <-s.done:
				slog.Info("Scheduler stopping")
				return
			case <-s.ticker.C:
				slog.Debug("Scheduler executing task")
				if err := task(); err != nil {
					slog.Error("Scheduler task error", "error", err)
				}
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	select {
	case s.done <- true:
	default:
	}
}
