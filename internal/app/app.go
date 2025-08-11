package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"dalnoboy/internal"
	"dalnoboy/internal/bot"
	"dalnoboy/internal/database"
	"dalnoboy/internal/service"
)

// App представляет основное приложение
type App struct {
	Name         string
	AdminBot     *bot.AdminBot
	DriverBot    *bot.DriverBot
	Database     *database.Database
	OrderService service.OrderService
	HTTPServer   *http.Server
}

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	AppName   string    `json:"app_name"`
}

// New создает новый экземпляр приложения
func New(name string) *App {
	return &App{
		Name: name,
	}
}

// healthCheckHandler обрабатывает health check запросы
func (a *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		AppName:   a.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getOrdersHandler обрабатывает запросы на получение всех заказов
func (a *App) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	orders, err := a.OrderService.GetAllOrders()
	if err != nil {
		log.Printf("Ошибка получения заказов: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// startHTTPServer запускает HTTP сервер
func (a *App) startHTTPServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.healthCheckHandler)
	mux.HandleFunc("/orders", a.getOrdersHandler)

	a.HTTPServer = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("🌐 HTTP сервер запущен на порту 8080")
	return a.HTTPServer.ListenAndServe()
}

// Run запускает приложение
func (a *App) Run() error {

	fmt.Printf("Приложение %s запущено\n", a.Name)

	// Создание конфига
	config, err := internal.NewConfig()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфига: %v", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации конфига: %v", err)
	}

	// Подключение к базе данных
	db, err := database.New(config)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}
	a.Database = db
	defer a.Database.Close()

	// Инициализация сервиса заказов
	a.OrderService = service.NewOrderService(db)

	// Инициализация админского бота
	adminBot, err := bot.NewAdminBot(config, db)
	if err != nil {
		return fmt.Errorf("ошибка инициализации админского бота: %v", err)
	}
	a.AdminBot = adminBot

	// Инициализация бота для водителей
	driverBot, err := bot.NewDriverBot(config, db)
	if err != nil {
		return fmt.Errorf("ошибка инициализации бота для водителей: %v", err)
	}
	a.DriverBot = driverBot

	// Запуск ботов и HTTP сервера в отдельных горутинах
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.startHTTPServer(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка HTTP сервера: %v", err)
		}
	}()

	fmt.Println("Оба бота и HTTP сервер запущены и работают...")
	wg.Wait()

	return nil
}

// Shutdown gracefully завершает работу приложения
func (a *App) Shutdown(ctx context.Context) error {
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Shutdown(ctx); err != nil {
			log.Printf("Ошибка завершения HTTP сервера: %v", err)
		}
	}
	return nil
}
