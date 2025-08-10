package bot

import (
	"fmt"
	"log"
	"strings"

	"dalnoboy/internal"
	"dalnoboy/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// mainMenuKeyboard возвращает главное меню с одной кнопкой "Заказы"
func mainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "📋 Заказы"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// ordersMenuKeyboard возвращает меню для раздела заказов
func ordersMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "⚙️ Фильтр"},
				{Text: "⬅️ Назад"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// filterMenuKeyboard возвращает меню настройки фильтров
func filterMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "📍 Маршрут"},
				{Text: "💰 Цена"},
			},
			{
				{Text: "📅 Дата"},
				{Text: "📦 Тип груза"},
			},
			{
				{Text: "♻️ Сбросить"},
				{Text: "⬅️ Назад"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// DriverBot представляет бота для водителей
type DriverBot struct {
	bot      *tgbotapi.BotAPI
	database *database.Database
}

// NewDriverBot создает новый экземпляр бота для водителей
func NewDriverBot(config *internal.Config, db *database.Database) (*DriverBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Bot.DriverToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Бот для водителей %s запущен", bot.Self.UserName)

	return &DriverBot{
		bot:      bot,
		database: db,
	}, nil
}

// Start запускает бота для водителей
func (db *DriverBot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := db.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			// Обработка сообщений
			db.handleMessage(update.Message)
		}
	}

	return nil
}

// formatOrders форматирует заказы для отображения без ID
func (db *DriverBot) formatOrders(orders []database.Order) string {
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

		result.WriteString(fmt.Sprintf("%d. 🚚 Заказ\n", i+1))
		result.WriteString(fmt.Sprintf("   📝 %s\n", order.Title))
		if order.Description != "" {
			result.WriteString(fmt.Sprintf("   📄 %s\n", order.Description))
		}
		result.WriteString(fmt.Sprintf("   📍 %s → %s\n", fromLoc, toLoc))
		result.WriteString(fmt.Sprintf("   ⚖️ %.1f кг\n", order.WeightKg))
		result.WriteString(fmt.Sprintf("   📏 %s\n", dimensions))
		result.WriteString(fmt.Sprintf("   🏷️ %s\n", tagsStr))
		result.WriteString(fmt.Sprintf("   💰 %.0f ₽\n", order.Price))
		result.WriteString(fmt.Sprintf("   📅 %s\n", dateStr))
		result.WriteString("\n")
	}

	return result.String()
}

// handleMessage обрабатывает входящие сообщения
func (db *DriverBot) handleMessage(message *tgbotapi.Message) {
	text := message.Text
	chatID := message.Chat.ID

	var response string
	var keyboard tgbotapi.ReplyKeyboardMarkup

	switch text {
	case "/start":
		response = "Добро пожаловать! Вы водитель. Выберите действие."
		keyboard = mainMenuKeyboard()
	case "/help", "❓ Помощь":
		response = "Доступные команды:\n/start - Начать работу\n/help - Показать помощь\n/orders - Посмотреть заказы\n/profile - Мой профиль"
	case "/orders", "📋 Заказы":
		// Получаем заказы из базы данных
		orders, err := db.database.GetOrders()
		if err != nil {
			log.Printf("Ошибка получения заказов: %v", err)
			response = "❌ Ошибка получения заказов из базы данных"
		} else {
			response = db.formatOrders(orders)
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
	db.bot.Send(msg)
}
