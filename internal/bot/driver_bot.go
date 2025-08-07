package bot

import (
	"log"

	"dalnoboy/internal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DriverBot представляет бота для водителей
type DriverBot struct {
	bot *tgbotapi.BotAPI
}

// NewDriverBot создает новый экземпляр бота для водителей
func NewDriverBot(config *internal.Config) (*DriverBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.DriverBotToken)
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
		if update.Message == nil {
			continue
		}

		// Обработка сообщений
		db.handleMessage(update.Message)
	}

	return nil
}

// handleMessage обрабатывает входящие сообщения
func (db *DriverBot) handleMessage(message *tgbotapi.Message) {
	text := message.Text
	chatID := message.Chat.ID

	var response string

	switch text {
	case "/start":
		response = "Добро пожаловать! Вы водитель. Используйте /help для списка команд."
	case "/help":
		response = "Доступные команды:\n/start - Начать работу\n/help - Показать помощь\n/orders - Посмотреть заказы\n/profile - Мой профиль"
	case "/orders":
		response = "Список доступных заказов:\n1. Заказ #123 - Москва → Санкт-Петербург\n2. Заказ #124 - Екатеринбург → Новосибирск"
	case "/profile":
		response = "Ваш профиль:\nИмя: Водитель\nРейтинг: 4.8\nЗаказов выполнено: 156"
	default:
		response = "Неизвестная команда. Используйте /help для списка команд."
	}

	msg := tgbotapi.NewMessage(chatID, response)
	db.bot.Send(msg)
}
