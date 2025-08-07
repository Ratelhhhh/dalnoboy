package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// BotConfig представляет конфигурацию ботов
type BotConfig struct {
	AdminToken  string
	DriverToken string
}

// DatabaseConfig представляет конфигурацию базы данных
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Config представляет общую конфигурацию приложения
type Config struct {
	Bot      BotConfig
	Database DatabaseConfig `yaml:"database"`
}

// NewConfig создает новый экземпляр конфига из YAML файла и переменных окружения
func NewConfig() (*Config, error) {
	configPath := "config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурационного файла: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %v", err)
	}

	// Загружаем токены ботов из переменных окружения
	config.Bot.AdminToken = os.Getenv("ADMIN_BOT_TOKEN")
	config.Bot.DriverToken = os.Getenv("DRIVER_BOT_TOKEN")

	return &config, nil
}

// Validate проверяет корректность конфига
func (c *Config) Validate() error {
	if c.Bot.AdminToken == "" {
		return &ConfigError{Field: "ADMIN_BOT_TOKEN", Message: "токен админского бота не установлен"}
	}
	if c.Bot.DriverToken == "" {
		return &ConfigError{Field: "DRIVER_BOT_TOKEN", Message: "токен бота водителя не установлен"}
	}
	if c.Database.Host == "" {
		return &ConfigError{Field: "database.host", Message: "хост базы данных не установлен"}
	}
	if c.Database.Port == 0 {
		return &ConfigError{Field: "database.port", Message: "порт базы данных не установлен"}
	}
	if c.Database.Name == "" {
		return &ConfigError{Field: "database.name", Message: "имя базы данных не установлено"}
	}
	if c.Database.User == "" {
		return &ConfigError{Field: "database.user", Message: "пользователь базы данных не установлен"}
	}
	if c.Database.Password == "" {
		return &ConfigError{Field: "database.password", Message: "пароль базы данных не установлен"}
	}
	return nil
}

// GetDBConnectionString возвращает строку подключения к базе данных
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)
}

// ConfigError представляет ошибку конфигурации
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}
