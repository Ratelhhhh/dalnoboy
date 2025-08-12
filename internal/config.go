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

// RedisConfig представляет конфигурацию Redis
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Config представляет общую конфигурацию приложения
type Config struct {
	Bot      BotConfig
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

// NewConfig создает новый экземпляр конфига из YAML файла и переменных окружения
func NewConfig() (*Config, error) {
	// Определяем путь к конфигурационному файлу
	configPath := "config.yaml"

	// Если указан путь через переменную окружения, используем его
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	} else {
		// Автоматически определяем окружение
		if isLocalDevelopment() {
			configPath = "config.local.yaml"
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурационного файла %s: %v", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %v", err)
	}

	// Загружаем токены ботов из переменных окружения
	config.Bot.AdminToken = os.Getenv("ADMIN_BOT_TOKEN")
	config.Bot.DriverToken = os.Getenv("DRIVER_BOT_TOKEN")

	// Загружаем настройки Redis из переменных окружения (приоритет над файлом)
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		config.Redis.Host = redisHost
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		if port, err := fmt.Sscanf(redisPort, "%d", &config.Redis.Port); err != nil || port == 0 {
			config.Redis.Port = 6379 // значение по умолчанию
		}
	}

	// Загружаем настройки базы данных из переменных окружения (приоритет над файлом)
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if port, err := fmt.Sscanf(dbPort, "%d", &config.Database.Port); err != nil || port == 0 {
			config.Database.Port = 5432 // значение по умолчанию
		}
	}

	return &config, nil
}

// isLocalDevelopment определяет, запущено ли приложение в локальной среде разработки
func isLocalDevelopment() bool {
	// Проверяем наличие переменной окружения
	if os.Getenv("ENV") == "local" || os.Getenv("ENV") == "development" {
		return true
	}

	// Проверяем, запущено ли приложение в Docker
	if os.Getenv("DOCKER_ENV") == "true" {
		return false
	}

	// Проверяем, существует ли локальный конфиг
	if _, err := os.Stat("config.local.yaml"); err == nil {
		return true
	}

	// По умолчанию считаем, что это локальная разработка
	return true
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
	if c.Redis.Host == "" {
		return &ConfigError{Field: "redis.host", Message: "хост Redis не установлен"}
	}
	if c.Redis.Port == 0 {
		return &ConfigError{Field: "redis.port", Message: "порт Redis не установлен"}
	}
	return nil
}

// GetDBConnectionString возвращает строку подключения к базе данных
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)
}

// GetRedisConnectionString возвращает строку подключения к Redis
func (c *Config) GetRedisConnectionString() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// ConfigError представляет ошибку конфигурации
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}
