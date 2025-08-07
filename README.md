# Dalnoboy

Telegram бот приложение для управления доставкой.

## Конфигурация

Приложение использует переменные окружения для конфигурации:

- `ADMIN_BOT_TOKEN` - токен админского бота
- `DRIVER_BOT_TOKEN` - токен бота для водителей

## Запуск

1. Установите переменные окружения:
```bash
export ADMIN_BOT_TOKEN="your_admin_bot_token_here"
export DRIVER_BOT_TOKEN="your_driver_bot_token_here"
```

2. Запустите приложение:
```bash
go run cmd/main.go
```

## Сборка

```bash
go build -o dalnoboy cmd/main.go
```

## Запуск собранного приложения

```bash
./dalnoboy
```

## Структура проекта

- `internal/config.go` - конфигурация приложения
- `internal/app/` - основная логика приложения
- `internal/bot/` - Telegram боты (админский и для водителей)
- `cmd/` - точка входа в приложение
