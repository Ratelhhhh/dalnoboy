package bot

import (
	"fmt"
	"log"
	"strings"

	"dalnoboy/internal"
	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"
	"dalnoboy/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AdminBot представляет админского бота
type AdminBot struct {
	bot          *tgbotapi.BotAPI
	database     *database.Database
	orderService service.OrderService
}

// NewAdminBot создает новый экземпляр админского бота
func NewAdminBot(config *internal.Config, db *database.Database) (*AdminBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Bot.AdminToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Админский бот %s запущен", bot.Self.UserName)

	return &AdminBot{
		bot:          bot,
		database:     db,
		orderService: service.NewOrderService(db),
	}, nil
}

// Start запускает админского бота
func (ab *AdminBot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := ab.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			// Обработка сообщений
			ab.handleMessage(update.Message)
		}
	}

	return nil
}

// formatOrders форматирует заказы для отображения с ID
func (ab *AdminBot) formatOrders(orders []domain.Order) string {
	if len(orders) == 0 {
		return "📋 Заказов пока нет"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("📋 Список доступных заказов (%d):\n\n", len(orders)))

	for i, order := range orders {
		// Форматируем дату
		dateStr := "Не указана"
		if order.AvailableFrom != nil {
			dateStr = order.AvailableFrom.Format("02.01.2006")
		}

		// Форматируем размеры
		dimensions := "Не указаны"
		if order.LengthCm != nil && order.WidthCm != nil && order.HeightCm != nil {
			dimensions = fmt.Sprintf("%.0f×%.0f×%.0f см", *order.LengthCm, *order.WidthCm, *order.HeightCm)
		}

		// Форматируем локации
		fromLoc := "Не указано"
		toLoc := "Не указано"
		if order.FromLocation != nil {
			fromLoc = *order.FromLocation
		}
		if order.ToLocation != nil {
			toLoc = *order.ToLocation
		}

		// Форматируем теги
		tagsStr := "Нет тегов"
		if len(order.Tags) > 0 {
			tagsStr = strings.Join(order.Tags, ", ")
		}

		result.WriteString(fmt.Sprintf("%d. 🚚 Заказ #%s\n", i+1, order.UUID[:8]))
		result.WriteString(fmt.Sprintf("   📝 %s\n", order.Title))
		if order.Description != "" {
			result.WriteString(fmt.Sprintf("   📄 %s\n", order.Description))
		}
		result.WriteString(fmt.Sprintf("   👤 %s (%s)\n", order.CustomerName, order.CustomerPhone))
		result.WriteString(fmt.Sprintf("   📍 %s → %s\n", fromLoc, toLoc))
		result.WriteString(fmt.Sprintf("   ⚖️ %.1f кг\n", order.WeightKg))
		result.WriteString(fmt.Sprintf("   📏 %s\n", dimensions))
		result.WriteString(fmt.Sprintf("   🏷️ %s\n", tagsStr))
		result.WriteString(fmt.Sprintf("   💰 %.0f ₽\n", order.Price))
		result.WriteString(fmt.Sprintf("   📅 %s\n", dateStr))
		result.WriteString(fmt.Sprintf("   🆔 ID: %s\n", order.UUID))
		result.WriteString("\n")
	}

	return result.String()
}

// handleMessage обрабатывает входящие сообщения
func (ab *AdminBot) handleMessage(message *tgbotapi.Message) {
	text := message.Text
	chatID := message.Chat.ID

	var response string
	var keyboard tgbotapi.ReplyKeyboardMarkup

	switch text {
	case "/start":
		response = "Добро пожаловать в админскую панель! Выберите действие."
		keyboard = mainMenuKeyboard()
	case "/help", "❓ Помощь":
		response = "Доступные команды:\n/start - Начать работу\n/help - Показать помощь\n/status - Статус системы\n/orders - Посмотреть заказы\n/filter - Настроить фильтры"
	case "/status":
		// Получаем статистику из базы данных
		ordersCount, err := ab.database.GetOrdersCount()
		if err != nil {
			log.Printf("Ошибка получения количества заказов: %v", err)
			ordersCount = -1
		}

		customersCount, err := ab.database.GetCustomersCount()
		if err != nil {
			log.Printf("Ошибка получения количества клиентов: %v", err)
			customersCount = -1
		}

		if ordersCount >= 0 && customersCount >= 0 {
			response = fmt.Sprintf("✅ Система работает нормально.\n📊 Статистика:\n📋 Заказов: %d\n👥 Клиентов: %d", ordersCount, customersCount)
		} else {
			response = "⚠️ Система работает, но есть проблемы с базой данных"
		}
	case "/orders", "📋 Заказы":
		// Получаем заказы через сервис
		orders, err := ab.orderService.GetAllOrders()
		if err != nil {
			log.Printf("Ошибка получения заказов: %v", err)
			response = "❌ Ошибка получения заказов из базы данных"
		} else {
			response = ab.formatOrders(orders)
		}
		keyboard = ordersMenuKeyboard()
	case "/filter", "⚙️ Фильтр":
		response = "Вы в меню фильтров. Выберите, что настроить:"
		keyboard = filterMenuKeyboard()
	case "📍 Маршрут":
		response = "Введите маршрут в сообщении, например: Москва → Санкт-Петербург"
		keyboard = filterMenuKeyboard()
	case "💰 Цена":
		response = "Укажите диапазон цены, например: 10000-20000"
		keyboard = filterMenuKeyboard()
	case "📅 Дата":
		response = "Укажите дату или диапазон, например: Сегодня или 2025-08-10 — 2025-08-15"
		keyboard = filterMenuKeyboard()
	case "📦 Тип груза":
		response = "Укажите тип груза, например: Рефрижератор, Негабарит, Опасный"
		keyboard = filterMenuKeyboard()
	case "♻️ Сбросить":
		response = "Фильтры сброшены"
		keyboard = filterMenuKeyboard()
	case "⬅️ Назад":
		response = "Главное меню"
		keyboard = mainMenuKeyboard()
	default:
		response = "Неизвестная команда. Используйте кнопки меню или /help для списка команд."
	}

	msg := tgbotapi.NewMessage(chatID, response)
	if keyboard.Keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	ab.bot.Send(msg)
}
