package domain

import (
	"time"

	"github.com/google/uuid"
)

// User представляет доменную модель пользователя
type User struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	TelegramID  *int64    `json:"telegram_id"`
	TelegramTag *string   `json:"telegram_tag"`
	CreatedAt   time.Time `json:"created_at"`
}
