#!/bin/bash

# Скрипт для запуска в Docker

echo "🐳 Запуск в Docker..."

# Проверяем, установлен ли Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен"
    exit 1
fi

# Проверяем, установлен ли Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен"
    exit 1
fi

echo "✅ Docker и Docker Compose установлены"

# Останавливаем существующие контейнеры
echo "🛑 Остановка существующих контейнеров..."
docker-compose down

# Запускаем все сервисы
echo "🚀 Запуск всех сервисов..."
docker-compose up -d

# Ждем, пока все сервисы будут готовы
echo "⏳ Ожидание готовности сервисов..."
sleep 10

# Проверяем статус
echo "📊 Статус сервисов:"
docker-compose ps

echo "✅ Приложение запущено в Docker!"
echo "🌐 Сайт: http://localhost:8080"
echo "🔌 API: http://localhost:8080/v1/orders"
echo "📱 Логи: docker-compose logs -f" 