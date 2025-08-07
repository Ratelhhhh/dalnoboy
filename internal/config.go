package internal

import (
	"os"
)

// Config представляет конфигурацию приложения
type Config struct {
	AdminBotToken  string
	DriverBotToken string
}

// NewConfig создает новый экземпляр конфига из переменных окружения
func NewConfig() *Config {
	return &Config{
		AdminBotToken:  os.Getenv("ADMIN_BOT_TOKEN"),
		DriverBotToken: os.Getenv("DRIVER_BOT_TOKEN"),
	}
}

// Validate проверяет корректность конфига
func (c *Config) Validate() error {
	if c.AdminBotToken == "" {
		return &ConfigError{Field: "ADMIN_BOT_TOKEN", Message: "токен админского бота не установлен"}
	}
	if c.DriverBotToken == "" {
		return &ConfigError{Field: "DRIVER_BOT_TOKEN", Message: "токен бота водителя не установлен"}
	}
	return nil
}

// ConfigError представляет ошибку конфигурации
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}
