package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-event-booker/internal/config"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"github.com/kstsm/wb-event-booker/internal/models"
	"net/http"
	"time"
)

type NotifierI interface {
	SendNotification(ctx context.Context, userID uuid.UUID, telegramID int64, message string) error
}

type TelegramNotifier struct {
	cfg    config.TelegramConfig
	client *http.Client
}

func NewTelegramNotifier(cfg config.TelegramConfig) NotifierI {
	return &TelegramNotifier{
		cfg: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (t *TelegramNotifier) SendNotification(
	ctx context.Context,
	userID uuid.UUID,
	telegramID int64,
	message string,
) error {

	if t.cfg.BotToken == "" {
		return fmt.Errorf("telegram bot token is not configured")
	}

	telegramMsg := models.TelegramMessage{
		ChatID: telegramID,
		Text:   message,
	}

	body, err := json.Marshal(telegramMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.cfg.BotToken)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	slog.Infof("Sending Telegram notification to user: user_id=%v, telegram_id=%d", userID, telegramID)
	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var telegramResp dto.TelegramResponse
		if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
			return fmt.Errorf("telegram API error: status %d", resp.StatusCode)
		}
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	var telegramResp dto.TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	return nil
}
