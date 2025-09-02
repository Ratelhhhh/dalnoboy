package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"dalnoboy/internal"
	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"
	"dalnoboy/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DriverBot представляет бота для водителей
type DriverBot struct {
	bot           *tgbotapi.BotAPI
	database      *database.Database
	orderService  *service.OrderService
	driverService *service.DriverService
}

// NewDriverBot создает новый экземпляр бота для водителей
func NewDriverBot(config *internal.Config, db *database.Database, orderService *service.OrderService, driverService *service.DriverService) (*DriverBot, error) {
	log.Printf("Инициализация бота для водителей с токеном: %s...", config.Bot.DriverToken[:10]+"...")

	bot, err := tgbotapi.NewBotAPI(config.Bot.DriverToken)
	if err != nil {
		log.Printf("Ошибка создания бота: %v", err)
		return nil, fmt.Errorf("ошибка создания бота для водителей: %v", err)
	}

	log.Printf("Бот для водителей %s запущен", bot.Self.UserName)

	return &DriverBot{
		bot:           bot,
		database:      db,
		orderService:  orderService,
		driverService: driverService,
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
func (db *DriverBot) formatOrders(orders []domain.Order) string {
	if len(orders) == 0 {
		return "📋 Заказов пока нет"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("📋 Список доступных заказов (%d):\n\n", len(orders)))

	for i, order := range orders {
		// Форматируем локации для межгородских перевозок
		fromLoc := "Не указано"
		toLoc := "Не указано"

		if order.FromCityName != nil && order.ToCityName != nil {
			// Основной маршрут между городами
			fromLoc = fmt.Sprintf("%s → %s", *order.FromCityName, *order.ToCityName)

			// Адреса в одной строке
			if order.FromAddress != nil && order.ToAddress != nil {
				toLoc = fmt.Sprintf("🏠 %s: %s | %s: %s",
					*order.FromCityName, *order.FromAddress,
					*order.ToCityName, *order.ToAddress)
			} else if order.FromAddress != nil {
				toLoc = fmt.Sprintf("🏠 %s: %s", *order.FromCityName, *order.FromAddress)
			} else if order.ToAddress != nil {
				toLoc = fmt.Sprintf("🏠 %s: %s", *order.ToCityName, *order.ToAddress)
			} else {
				toLoc = "🏠 Адреса не указаны"
			}
		} else if order.FromCityName != nil {
			fromLoc = *order.FromCityName
			if order.FromAddress != nil {
				toLoc = fmt.Sprintf("🏠 Адрес: %s", *order.FromAddress)
			}
		} else if order.ToCityName != nil {
			toLoc = *order.ToCityName
			if order.ToAddress != nil {
				fromLoc = fmt.Sprintf("🏠 Адрес: %s", *order.ToAddress)
			}
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
		result.WriteString(fmt.Sprintf("   %s\n", fromLoc))
		result.WriteString(fmt.Sprintf("   %s\n", toLoc))
		result.WriteString(fmt.Sprintf("   ⚖️ %.1f кг | 💰 %.0f ₽\n", order.WeightKg, order.Price))
		result.WriteString(fmt.Sprintf("   👤 %s | 📱 %s\n", order.CustomerName, order.CustomerPhone))

		// Добавляем теги только если они есть
		if len(order.Tags) > 0 {
			result.WriteString(fmt.Sprintf("   🏷️ %s\n", tagsStr))
		}

		result.WriteString("\n")
	}
	return result.String()
}

// splitMessage разбивает длинное сообщение на части для Telegram
func (db *DriverBot) splitMessage(text string, maxLength int) []string {
	if len(text) <= maxLength {
		return []string{text}
	}

	var parts []string
	var currentPart strings.Builder
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		// Если добавление текущей строки превысит лимит
		if currentPart.Len()+len(line)+1 > maxLength {
			// Сохраняем текущую часть
			if currentPart.Len() > 0 {
				parts = append(parts, strings.TrimSpace(currentPart.String()))
				currentPart.Reset()
			}

			// Если одна строка слишком длинная, разбиваем её
			if len(line) > maxLength {
				// Разбиваем длинную строку по словам
				words := strings.Fields(line)
				var tempLine strings.Builder

				for _, word := range words {
					if tempLine.Len()+len(word)+1 > maxLength {
						if tempLine.Len() > 0 {
							parts = append(parts, strings.TrimSpace(tempLine.String()))
							tempLine.Reset()
						}
					}
					if tempLine.Len() > 0 {
						tempLine.WriteString(" ")
					}
					tempLine.WriteString(word)
				}

				if tempLine.Len() > 0 {
					currentPart.WriteString(tempLine.String())
					currentPart.WriteString("\n")
				}
			} else {
				currentPart.WriteString(line)
				currentPart.WriteString("\n")
			}
		} else {
			currentPart.WriteString(line)
			currentPart.WriteString("\n")
		}
	}

	// Добавляем последнюю часть
	if currentPart.Len() > 0 {
		parts = append(parts, strings.TrimSpace(currentPart.String()))
	}

	return parts
}

// handleMessage обрабатывает входящие сообщения
func (db *DriverBot) handleMessage(message *tgbotapi.Message) {
	text := message.Text
	chatID := message.Chat.ID

	// Автоматическая регистрация водителя по первому контакту с ботом
	telegramID := message.From.ID
	var tag *string
	if message.From.UserName != "" {
		uname := "@" + message.From.UserName
		tag = &uname
	}
	name := strings.TrimSpace(message.From.FirstName + " " + message.From.LastName)
	if name == "" {
		name = message.From.UserName
	}
	driver, ensureErr := db.driverService.EnsureDriverExistsByTelegram(name, telegramID, tag)
	if ensureErr != nil {
		log.Printf("Не удалось авто-регистрировать водителя %d: %v", telegramID, ensureErr)
	}

	var response string
	var keyboard tgbotapi.ReplyKeyboardMarkup

	switch text {
	case "/start":
		response = "Добро пожаловать! Вы водитель. Выберите действие."
		keyboard = driverMainMenuKeyboard()
	case "/help", "❓ Помощь":
		response = "Доступные команды:\n/start - Начать работу\n/help - Показать помощь\n/orders - Посмотреть заказы\n🔔 Включить уведомления - Получать новые заказы\n🔕 Выключить уведомления - Отключить получение заказов"
	case "/orders", "📋 Заказы":
		// Получаем только активные заказы через сервис
		orders, err := db.orderService.GetActiveOrders()
		if err != nil {
			log.Printf("Ошибка получения активных заказов: %v", err)
			response = "❌ Ошибка получения заказов из базы данных"
		} else {
			response = db.formatOrders(orders)
		}
		keyboard = driverMainMenuKeyboard()
	case "🔔 Включить уведомления":
		// Включаем уведомления для текущего водителя
		if driver == nil {
			log.Printf("Водитель с Telegram ID %d не найден при включении уведомлений", telegramID)
			response = "❌ Не удалось включить уведомления: водитель не найден."
			break
		}
		if err := db.driverService.UpdateDriverNotifications(driver.UUID, true); err != nil {
			log.Printf("Ошибка обновления статуса уведомлений для водителя %s: %v", driver.UUID, err)
			response = "❌ Не удалось включить уведомления. Попробуйте позже."
		} else {
			response = "✅ Уведомления включены! Теперь вы будете получать новые заказы согласно вашим настройкам."
		}
		keyboard = driverMainMenuKeyboard()
	case "🔕 Выключить уведомления":
		// Выключаем уведомления для текущего водителя
		if driver == nil {
			log.Printf("Водитель с Telegram ID %d не найден при выключении уведомлений", telegramID)
			response = "❌ Не удалось выключить уведомления: водитель не найден."
			break
		}
		if err := db.driverService.UpdateDriverNotifications(driver.UUID, false); err != nil {
			log.Printf("Ошибка обновления статуса уведомлений для водителя %s: %v", driver.UUID, err)
			response = "❌ Не удалось выключить уведомления. Попробуйте позже."
		} else {
			response = "🔕 Уведомления выключены. Вы не будете получать новые заказы."
		}
		keyboard = driverMainMenuKeyboard()

	// Закомментировано - убираем функционал фильтров
	/*
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
	*/

	case "⬅️ Назад":
		response = "Главное меню"
		keyboard = driverMainMenuKeyboard()
	default:
		response = "Неизвестная команда. Используйте кнопки меню или /help для списка команд."
	}

	msg := tgbotapi.NewMessage(chatID, response)
	if keyboard.Keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	// Разбиваем длинные сообщения на части
	parts := db.splitMessage(response, 4000) // Telegram лимит 4096, берем с запасом

	for i, part := range parts {
		// Для первой части используем клавиатуру, для остальных - нет
		if i == 0 {
			msg := tgbotapi.NewMessage(chatID, part)
			if keyboard.Keyboard != nil {
				msg.ReplyMarkup = keyboard
			}
			if _, err := db.bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки части сообщения %d: %v", i+1, err)
			}
		} else {
			// Для дополнительных частей без клавиатуры
			msg := tgbotapi.NewMessage(chatID, part)
			if _, err := db.bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки части сообщения %d: %v", i+1, err)
			}
		}

		// Небольшая задержка между сообщениями
		if i < len(parts)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
