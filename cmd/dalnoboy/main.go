package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dalnoboy/internal/app"
)

func main() {
	appName := "dalnoboy"
	if len(os.Args) > 1 {
		appName = os.Args[1]
	}

	// Устанавливаем локальную конфигурацию по умолчанию
	if os.Getenv("CONFIG_PATH") == "" {
		os.Setenv("CONFIG_PATH", "config.local.yaml")
	}

	application := app.New(appName)

	// Канал для сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск приложения в горутине
	go func() {
		if err := application.Run(); err != nil {
			log.Printf("Ошибка приложения: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Получен сигнал завершения, завершаю работу...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении: %v", err)
	}

	log.Println("Приложение завершено")
}
