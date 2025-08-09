package bot

import (
	"log"

	"dalnoboy/internal"

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
	bot *tgbotapi.BotAPI
}

// NewDriverBot создает новый экземпляр бота для водителей
func NewDriverBot(config *internal.Config) (*DriverBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Bot.DriverToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Бот для водителей %s запущен", bot.Self.UserName)

	return &DriverBot{
		bot: bot,
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
		response = "📋 Список доступных заказов:\n\n1. 🚚 Заказ #123\n   📍 Москва → Санкт-Петербург\n   💰 15,000 ₽\n   📅 Сегодня\n\n2. 🚚 Заказ #124\n   📍 Екатеринбург → Новосибирск\n   💰 12,000 ₽\n   📅 Завтра"
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
