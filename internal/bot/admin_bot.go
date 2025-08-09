package bot

import (
	"log"

	"dalnoboy/internal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AdminBot представляет админского бота
type AdminBot struct {
	bot *tgbotapi.BotAPI
}

// NewAdminBot создает новый экземпляр админского бота
func NewAdminBot(config *internal.Config) (*AdminBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Bot.AdminToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Админский бот %s запущен", bot.Self.UserName)

	return &AdminBot{
		bot: bot,
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
		response = "Система работает нормально. Все боты активны."
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
	ab.bot.Send(msg)
}
