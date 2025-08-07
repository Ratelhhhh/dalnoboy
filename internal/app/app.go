package app

import (
	"fmt"
	"log"
	"sync"

	"dalnoboy/internal"
	"dalnoboy/internal/bot"
)

// App представляет основное приложение
type App struct {
	Name      string
	AdminBot  *bot.AdminBot
	DriverBot *bot.DriverBot
}

// New создает новый экземпляр приложения
func New(name string) *App {
	return &App{
		Name: name,
	}
}

// Run запускает приложение
func (a *App) Run() error {
	fmt.Printf("Приложение %s запущено\n", a.Name)

	// Создание конфига
	config := internal.NewConfig()
	if err := config.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации конфига: %v", err)
	}

	// Инициализация админского бота
	adminBot, err := bot.NewAdminBot(config)
	if err != nil {
		return fmt.Errorf("ошибка инициализации админского бота: %v", err)
	}
	a.AdminBot = adminBot

	// Инициализация бота для водителей
	driverBot, err := bot.NewDriverBot(config)
	if err != nil {
		return fmt.Errorf("ошибка инициализации бота для водителей: %v", err)
	}
	a.DriverBot = driverBot

	// Запуск ботов в отдельных горутинах
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.AdminBot.Start(); err != nil {
			log.Printf("Ошибка админского бота: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.DriverBot.Start(); err != nil {
			log.Printf("Ошибка бота для водителей: %v", err)
		}
	}()

	fmt.Println("Оба бота запущены и работают...")
	wg.Wait()

	return nil
}
