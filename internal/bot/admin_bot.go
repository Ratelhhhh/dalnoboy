package bot

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AdminBot представляет админского бота
type AdminBot struct {
	bot *tgbotapi.BotAPI
}

// NewAdminBot создает новый экземпляр админского бота
func NewAdminBot() (*AdminBot, error) {
	token := os.Getenv("ADMIN_BOT_TOKEN")
	if token == "" {
		panic("ADMIN_BOT_TOKEN не установлен в переменных окружения")
	}

	bot, err := tgbotapi.NewBotAPI(token)
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
		if update.Message == nil {
			continue
		}

		// Обработка сообщений
		ab.handleMessage(update.Message)
	}

	return nil
}

// handleMessage обрабатывает входящие сообщения
func (ab *AdminBot) handleMessage(message *tgbotapi.Message) {
	text := message.Text
	chatID := message.Chat.ID

	var response string

	switch text {
	case "/start":
		response = "Добро пожаловать в админскую панель! Используйте /help для списка команд."
	case "/help":
		response = "Доступные команды:\n/start - Начать работу\n/help - Показать помощь\n/status - Статус системы"
	case "/status":
		response = "Система работает нормально. Все боты активны."
	default:
		response = "Неизвестная команда. Используйте /help для списка команд."
	}

	msg := tgbotapi.NewMessage(chatID, response)
	ab.bot.Send(msg)
}
