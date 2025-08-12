#!/bin/bash

# Скрипт для запуска локальной разработки

echo "🚀 Запуск локальной разработки..."

# Проверяем, запущен ли Redis (локально или в Docker)
if ! redis-cli ping > /dev/null 2>&1; then
    echo "🔍 Redis не найден локально, проверяем Docker..."
    
    # Проверяем, запущен ли Redis в Docker
    if docker ps | grep -q "dalnoboy-redis"; then
        echo "✅ Redis запущен в Docker"
        # Устанавливаем переменные для подключения к Docker Redis
        export REDIS_HOST=localhost
        export REDIS_PORT=6379
    else
        echo "❌ Redis не запущен ни локально, ни в Docker"
        echo "💡 Запустите Redis одним из способов:"
        echo "   1. sudo systemctl start redis-server"
        echo "   2. docker run -d -p 6379:6379 redis:7-alpine"
        echo "   3. docker-compose up -d redis"
        exit 1
    fi
else
    echo "✅ Redis запущен локально"
fi

# Проверяем, запущена ли PostgreSQL локально
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "🔍 PostgreSQL не найден локально, проверяем Docker..."
    
    # Проверяем, запущена ли PostgreSQL в Docker
    if docker ps | grep -q "dalnoboy-postgres"; then
        echo "✅ PostgreSQL запущен в Docker"
        # Устанавливаем переменные для подключения к Docker PostgreSQL
        export DB_HOST=localhost
        export DB_PORT=5432
    else
        echo "❌ PostgreSQL не запущен ни локально, ни в Docker"
        echo "💡 Запустите PostgreSQL одним из способов:"
        echo "   1. sudo systemctl start postgresql"
        echo "   2. docker run -d -p 5432:5432 -e POSTGRES_DB=dalnoboy -e POSTGRES_USER=dalnoboy -e POSTGRES_PASSWORD=dalnoboy_password postgres:15-alpine"
        echo "   3. docker-compose up -d postgres"
        exit 1
    fi
else
    echo "✅ PostgreSQL запущен локально"
fi

echo "✅ Все необходимые сервисы доступны"

# Устанавливаем переменные окружения для локальной разработки
export ENV=local
export CONFIG_PATH=config.local.yaml

echo "🔧 Используется конфигурация: $CONFIG_PATH"
echo "🌍 Окружение: $ENV"

# Запускаем приложение
echo "🚀 Запуск приложения..."
go run ./cmd/dalnoboy 