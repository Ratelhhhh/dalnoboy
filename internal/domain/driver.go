package domain

import (
	"time"

	"github.com/google/uuid"
)

// Driver представляет доменную модель водителя
type Driver struct {
	UUID                uuid.UUID  `json:"uuid"`
	Name                string     `json:"name"`
	TelegramID          int64      `json:"telegram_id"`
	TelegramTag         *string    `json:"telegram_tag"`
	NotificationEnabled bool       `json:"notification_enabled"`
	CityUUID            *uuid.UUID `json:"city_uuid"`
	CityName            *string    `json:"city_name"`
	CreatedAt           time.Time  `json:"created_at"`
}

// SetCityAndNotificationRequest представляет запрос на обновление города и уведомлений водителя
type SetCityAndNotificationRequest struct {
	DriverUUID          uuid.UUID `json:"driver_uuid"`
	CityName            string    `json:"city_name"`
	NotificationEnabled *bool     `json:"notification_enabled"`
}
