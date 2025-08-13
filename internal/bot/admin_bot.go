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
	userService  *service.UserService
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
		userService:  service.NewUserService(db),
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

// formatUsers форматирует пользователей для отображения
func (ab *AdminBot) formatUsers(users []domain.User) string {
	if len(users) == 0 {
		return "👥 Пользователей пока нет"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("👥 Список пользователей (%d):\n\n", len(users)))

	for i, user := range users {
		// Форматируем Telegram ID
		telegramIDStr := "-"
		if user.TelegramID != nil {
			telegramIDStr = fmt.Sprintf("%d", *user.TelegramID)
		}

		// Форматируем Telegram Tag
		telegramTagStr := "-"
		if user.TelegramTag != nil {
			telegramTagStr = *user.TelegramTag
		}

		result.WriteString(fmt.Sprintf("%d. 👤 %s\n", i+1, user.Name))
		result.WriteString(fmt.Sprintf("   📱 %s\n", user.Phone))
		result.WriteString(fmt.Sprintf("   🆔 Telegram ID: %s\n", telegramIDStr))
		result.WriteString(fmt.Sprintf("   🏷️ Telegram Tag: %s\n", telegramTagStr))
		result.WriteString(fmt.Sprintf("   📅 Создан: %s\n", user.CreatedAt.Format("02.01.2006 15:04")))
		result.WriteString(fmt.Sprintf("   🆔 UUID: %s\n", user.UUID.String()))
		result.WriteString("\n")
	}

	return result.String()
}

// parseUserMessage парсит сообщение с данными пользователя
// Формат: ADD_USER\nИмя\nТелефон\nTelegramID\nTelegramTag
func (ab *AdminBot) parseUserMessage(text string) (*domain.User, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("недостаточно данных. Нужно минимум: имя, телефон")
	}

	// Проверяем ключ
	if lines[0] != "ADD_USER" {
		return nil, fmt.Errorf("неверный ключ. Ожидается ADD_USER")
	}

	name := strings.TrimSpace(lines[1])
	phone := strings.TrimSpace(lines[2])

	if name == "" || phone == "" {
		return nil, fmt.Errorf("имя и телефон не могут быть пустыми")
	}

	var telegramID *int64
	var telegramTag *string

	// Telegram ID (если указан)
	if len(lines) > 3 && strings.TrimSpace(lines[3]) != "-" && strings.TrimSpace(lines[3]) != "" {
		if id, err := parseTelegramID(lines[3]); err == nil {
			telegramID = &id
		}
	}

	// Telegram Tag (если указан)
	if len(lines) > 4 && strings.TrimSpace(lines[4]) != "-" && strings.TrimSpace(lines[4]) != "" {
		tag := strings.TrimSpace(lines[4])
		if tag != "" {
			telegramTag = &tag
		}
	}

	user := &domain.User{
		Name:        name,
		Phone:       phone,
		TelegramID:  telegramID,
		TelegramTag: telegramTag,
	}

	return user, nil
}

// parseTelegramID парсит Telegram ID из строки
func parseTelegramID(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &id)
	return id, err
}

// formatTelegramID форматирует Telegram ID для отображения
func formatTelegramID(id *int64) string {
	if id == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *id)
}

// formatTelegramTag форматирует Telegram Tag для отображения
func formatTelegramTag(tag *string) string {
	if tag == nil {
		return "-"
	}
	return *tag
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
		keyboard = adminMainMenuKeyboard()
	case "/help", "❓ Помощь":
		response = "Доступные команды:\n/start - Начать работу\n/help - Показать помощь\n/status - Статус системы\n/orders - Посмотреть заказы\n/users - Посмотреть пользователей\n/filter - Настроить фильтры"
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
	case "/users", "👥 Пользователи":
		// Получаем пользователей через сервис
		users, err := ab.userService.GetAllUsers()
		if err != nil {
			log.Printf("Ошибка получения пользователей: %v", err)
			response = "❌ Ошибка получения пользователей из базы данных"
		} else {
			response = ab.formatUsers(users)
		}
		keyboard = usersMenuKeyboard()
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
		keyboard = adminMainMenuKeyboard()
	default:
		// Проверяем, не является ли сообщение данными для добавления пользователя
		if strings.HasPrefix(text, "ADD_USER") {
			user, err := ab.parseUserMessage(text)
			if err != nil {
				response = fmt.Sprintf("❌ Ошибка парсинга данных пользователя: %v\n\nПример правильного формата:\nADD_USER\nИван Иванов\n+79001234567\n123456789\n@ivan_username", err)
			} else {
				// Создаем пользователя через сервис
				createdUser, err := ab.userService.CreateUser(user.Name, user.Phone, user.TelegramID, user.TelegramTag)
				if err != nil {
					response = fmt.Sprintf("❌ Ошибка создания пользователя: %v", err)
				} else {
					response = fmt.Sprintf("✅ Пользователь успешно создан!\n\n👤 Имя: %s\n📱 Телефон: %s\n🆔 Telegram ID: %s\n🏷️ Telegram Tag: %s\n🆔 UUID: %s",
						createdUser.Name,
						createdUser.Phone,
						formatTelegramID(createdUser.TelegramID),
						formatTelegramTag(createdUser.TelegramTag),
						createdUser.UUID)
				}
			}
		} else {
			response = "Неизвестная команда. Используйте кнопки меню или /help для списка команд.\n\nДля добавления пользователя используйте формат:\nADD_USER\nИмя\nТелефон\nTelegramID\nTelegramTag"
		}
	}

	msg := tgbotapi.NewMessage(chatID, response)
	if keyboard.Keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	ab.bot.Send(msg)
}
