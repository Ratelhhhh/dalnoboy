package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"dalnoboy/internal"
	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"
	"dalnoboy/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

// AdminBot представляет админского бота
type AdminBot struct {
	bot             *tgbotapi.BotAPI
	database        *database.Database
	orderService    *service.OrderService
	customerService *service.CustomerService
	driverService   *service.DriverService
}

// NewAdminBot создает новый экземпляр админского бота
func NewAdminBot(config *internal.Config, db *database.Database, orderService *service.OrderService, customerService *service.CustomerService, driverService *service.DriverService) (*AdminBot, error) {
	log.Printf("Инициализация админского бота с токеном: %s...", config.Bot.AdminToken[:10]+"...")

	bot, err := tgbotapi.NewBotAPI(config.Bot.AdminToken)
	if err != nil {
		log.Printf("Ошибка создания бота: %v", err)
		return nil, fmt.Errorf("ошибка создания админского бота: %v", err)
	}

	log.Printf("Админский бот %s запущен", bot.Self.UserName)

	return &AdminBot{
		bot:             bot,
		database:        db,
		orderService:    orderService,
		customerService: customerService,
		driverService:   driverService,
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
	result.WriteString(fmt.Sprintf("📋 Список заказов (%d):\n\n", len(orders)))

	for i, order := range orders {
		// Форматируем дату
		dateStr := ""
		if order.AvailableFrom != nil {
			dateStr = order.AvailableFrom.Format("02.01.2006")
		}

		// Форматируем размеры
		dimensions := "Не указаны"
		if order.LengthCm != nil && order.WidthCm != nil && order.HeightCm != nil {
			dimensions = fmt.Sprintf("%.0f×%.0f×%.0f см", *order.LengthCm, *order.WidthCm, *order.HeightCm)
		}

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

		// Форматируем статус
		statusEmoji := "🟢"
		statusText := "Активный"
		if order.Status == "archived" {
			statusEmoji = "🔴"
			statusText = "Архивный"
		}

		result.WriteString(fmt.Sprintf("%d. 🚚 Заказ #%s\n", i+1, order.UUID[:8]))
		result.WriteString(fmt.Sprintf("   %s %s\n", statusEmoji, statusText))
		result.WriteString(fmt.Sprintf("   📝 %s\n", order.Title))
		if order.Description != "" {
			result.WriteString(fmt.Sprintf("   📄 %s\n", order.Description))
		}
		result.WriteString(fmt.Sprintf("   👤 %s (%s)\n", order.CustomerName, order.CustomerPhone))

		// Добавляем информацию о Telegram заказчика
		if order.CustomerTelegramID != nil {
			result.WriteString(fmt.Sprintf("   🆔 Telegram ID: %d\n", *order.CustomerTelegramID))
		}
		if order.CustomerTelegramTag != nil && *order.CustomerTelegramTag != "" {
			result.WriteString(fmt.Sprintf("   🏷️ Telegram: %s\n", *order.CustomerTelegramTag))
		}
		result.WriteString(fmt.Sprintf("   %s\n", fromLoc))
		result.WriteString(fmt.Sprintf("   %s\n", toLoc))
		result.WriteString(fmt.Sprintf("   ⚖️ %.1f кг\n", order.WeightKg))
		result.WriteString(fmt.Sprintf("   📏 %s\n", dimensions))
		result.WriteString(fmt.Sprintf("   🏷️ %s\n", tagsStr))
		result.WriteString(fmt.Sprintf("   💰 %.0f ₽\n", order.Price))
		if dateStr != "" {
			result.WriteString(fmt.Sprintf("   📅 %s\n", dateStr))
		}
		result.WriteString(fmt.Sprintf("   🆔 ID: %s\n", order.UUID))
		result.WriteString("\n")
	}

	return result.String()
}

// formatCustomers форматирует заказчиков для отображения
func (ab *AdminBot) formatCustomers(customers []domain.Customer) string {
	if len(customers) == 0 {
		return "👥 Заказчиков пока нет"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("👥 Список заказчиков (%d):\n\n", len(customers)))

	for i, customer := range customers {
		// Форматируем Telegram ID
		telegramIDStr := "-"
		if customer.TelegramID != nil {
			telegramIDStr = fmt.Sprintf("%d", *customer.TelegramID)
		}

		// Форматируем Telegram Tag
		telegramTagStr := "-"
		if customer.TelegramTag != nil {
			telegramTagStr = *customer.TelegramTag
		}

		result.WriteString(fmt.Sprintf("%d. 👤 %s\n", i+1, customer.Name))
		result.WriteString(fmt.Sprintf("   📱 %s\n", customer.Phone))
		result.WriteString(fmt.Sprintf("   🆔 Telegram ID: %s\n", telegramIDStr))
		result.WriteString(fmt.Sprintf("   🏷️ Telegram Tag: %s\n", telegramTagStr))
		result.WriteString(fmt.Sprintf("   📅 Создан: %s\n", customer.CreatedAt.Format("02.01.2006 15:04")))
		result.WriteString(fmt.Sprintf("   🆔 UUID: %s\n", customer.UUID.String()))
		result.WriteString("\n")
	}

	return result.String()
}

// formatDrivers форматирует водителей для отображения
func (ab *AdminBot) formatDrivers(drivers []domain.Driver) string {
	if len(drivers) == 0 {
		return "🚚 Водителей пока нет"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("🚚 Список водителей (%d):\n\n", len(drivers)))

	for i, driver := range drivers {
		// Форматируем Telegram Tag
		telegramTagStr := "-"
		if driver.TelegramTag != nil {
			telegramTagStr = *driver.TelegramTag
		}

		// Форматируем название города
		cityNameStr := "Не указан"
		if driver.CityName != nil {
			cityNameStr = *driver.CityName
		}

		// Форматируем статус уведомлений
		notificationStatus := "🔔 Включены"
		if !driver.NotificationEnabled {
			notificationStatus = "🔕 Выключены"
		}

		result.WriteString(fmt.Sprintf("%d. 🚚 %s\n", i+1, driver.Name))
		result.WriteString(fmt.Sprintf("   📱 Telegram ID: %d\n", driver.TelegramID))
		result.WriteString(fmt.Sprintf("   🏷️ Telegram Tag: %s\n", telegramTagStr))
		result.WriteString(fmt.Sprintf("   🏙️ Город: %s\n", cityNameStr))
		result.WriteString(fmt.Sprintf("   %s\n", notificationStatus))
		result.WriteString(fmt.Sprintf("   📅 Зарегистрирован: %s\n", driver.CreatedAt.Format("02.01.2006 15:04")))
		result.WriteString(fmt.Sprintf("   🆔 UUID: %s\n", driver.UUID))
		result.WriteString("\n")
	}

	return result.String()
}

// parseCustomerMessage парсит сообщение с данными заказчика
// Формат: ADD_USER\nИмя\nТелефон\nTelegramID\nTelegramTag
func (ab *AdminBot) parseCustomerMessage(text string) (*domain.Customer, error) {
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

	customer := &domain.Customer{
		Name:        name,
		Phone:       phone,
		TelegramID:  telegramID,
		TelegramTag: telegramTag,
	}

	return customer, nil
}

// parseOrderMessage парсит сообщение с данными заказа (старый формат)
// Формат: ADD_ORDER\nНазвание\nОписание\nВес\nДлина\nШирина\nВысота\nОткуда город UUID\nОткуда адрес\nКуда город UUID\nКуда адрес\nТеги\nЦена\nДата\nUUID клиента
func (ab *AdminBot) parseOrderMessage(text string) (*domain.Order, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 15 {
		return nil, fmt.Errorf("недостаточно данных. Нужно 15 строк: ADD_ORDER + 14 полей (старый формат)")
	}

	// Проверяем ключ
	if lines[0] != "ADD_ORDER" {
		return nil, fmt.Errorf("неверный ключ. Ожидается ADD_ORDER")
	}

	title := strings.TrimSpace(lines[1])
	description := strings.TrimSpace(lines[2])

	if title == "" || description == "" {
		return nil, fmt.Errorf("название и описание не могут быть пустыми")
	}

	// Парсим вес
	weightKg, err := parseFloat(lines[3])
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга веса: %v", err)
	}

	// Парсим размеры (опционально)
	lengthCm := parseOptionalFloat(lines[4])
	widthCm := parseOptionalFloat(lines[5])
	heightCm := parseOptionalFloat(lines[6])

	// Парсим локации (опционально)
	fromCityUUID := parseOptionalString(lines[7])
	fromAddress := parseOptionalString(lines[8])
	toCityUUID := parseOptionalString(lines[9])
	toAddress := parseOptionalString(lines[10])

	// Парсим теги
	var tags []string
	if lines[11] != "-" && strings.TrimSpace(lines[11]) != "" {
		tags = strings.Split(strings.TrimSpace(lines[11]), ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// Парсим цену
	price, err := parseFloat(lines[12])
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга цены: %v", err)
	}

	// Парсим дату (опционально)
	availableFrom := parseOptionalDate(lines[13])

	// Парсим UUID клиента
	customerUUID := strings.TrimSpace(lines[14])
	if customerUUID == "" {
		return nil, fmt.Errorf("UUID клиента не может быть пустым")
	}

	order := &domain.Order{
		Title:         title,
		Description:   description,
		WeightKg:      weightKg,
		LengthCm:      lengthCm,
		WidthCm:       widthCm,
		HeightCm:      heightCm,
		FromCityUUID:  fromCityUUID,
		FromAddress:   fromAddress,
		ToCityUUID:    toCityUUID,
		ToAddress:     toAddress,
		Tags:          tags,
		Price:         price,
		AvailableFrom: availableFrom,
		CustomerUUID:  customerUUID,
	}

	return order, nil
}

// parseSetCityAndNotificationMessage парсит сообщение с командой настройки города и уведомлений водителя
// Формат: SET_CITY_AND_NOTIFICATION, <driver_uuid>, <city_name_or_-_or_empty>, <вкл/выкл_or_empty>
func (ab *AdminBot) parseSetCityAndNotificationMessage(text string) (*domain.SetCityAndNotificationRequest, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("недостаточно данных. Нужно минимум: SET_CITY_AND_NOTIFICATION + параметры")
	}

	// Проверяем ключ
	if lines[0] != "SET_CITY_AND_NOTIFICATION" {
		return nil, fmt.Errorf("неверный ключ. Ожидается SET_CITY_AND_NOTIFICATION")
	}

	// Парсим параметры
	params := strings.Split(lines[1], ",")
	if len(params) < 1 {
		return nil, fmt.Errorf("неверный формат параметров. Ожидается: UUID, город, уведомления")
	}

	// UUID водителя (обязательный параметр)
	driverUUIDStr := strings.TrimSpace(params[0])
	if driverUUIDStr == "" {
		return nil, fmt.Errorf("UUID водителя не может быть пустым")
	}

	driverUUID, err := uuid.Parse(driverUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат UUID водителя: %v", err)
	}

	// Город (опциональный параметр)
	var cityName string
	if len(params) > 1 {
		cityName = strings.TrimSpace(params[1])
	}

	// Статус уведомлений (опциональный параметр)
	var notificationEnabled *bool
	if len(params) > 2 {
		notificationStr := strings.TrimSpace(params[2])
		if notificationStr != "" {
			switch strings.ToLower(notificationStr) {
			case "вкл", "включить", "true", "1", "on":
				enabled := true
				notificationEnabled = &enabled
			case "выкл", "выключить", "false", "0", "off":
				enabled := false
				notificationEnabled = &enabled
			default:
				return nil, fmt.Errorf("неверный формат статуса уведомлений. Используйте: вкл/выкл, включить/выключить, true/false, 1/0, on/off")
			}
		}
	}

	request := &domain.SetCityAndNotificationRequest{
		DriverUUID:          driverUUID,
		CityName:            cityName,
		NotificationEnabled: notificationEnabled,
	}

	return request, nil
}

// parseOrderTgMessage парсит упрощенное сообщение с данными заказа для Telegram
// Формат: ADD_ORDER\nНазвание\nОписание\nВес\nОткуда город\nОткуда адрес\nКуда город\nКуда адрес\nЦена\nUUID клиента
func (ab *AdminBot) parseOrderTgMessage(text string) (*domain.CreateOrderTgRequest, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 10 {
		return nil, fmt.Errorf("недостаточно данных. Нужно 10 строк: ADD_ORDER + 9 полей")
	}

	// Проверяем ключ
	if lines[0] != "ADD_ORDER" {
		return nil, fmt.Errorf("неверный ключ. Ожидается ADD_ORDER")
	}

	title := strings.TrimSpace(lines[1])
	description := strings.TrimSpace(lines[2])

	if title == "" || description == "" {
		return nil, fmt.Errorf("название и описание не могут быть пустыми")
	}

	// Парсим вес
	weightKg, err := parseFloat(lines[3])
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга веса: %v", err)
	}

	// Парсим локации
	fromCityName := strings.TrimSpace(lines[4])
	fromAddress := strings.TrimSpace(lines[5])
	toCityName := strings.TrimSpace(lines[6])
	toAddress := strings.TrimSpace(lines[7])

	if fromCityName == "" || fromAddress == "" || toCityName == "" || toAddress == "" {
		return nil, fmt.Errorf("названия городов и адреса не могут быть пустыми")
	}

	// Парсим цену
	price, err := parseFloat(lines[8])
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга цены: %v", err)
	}

	// Парсим UUID клиента
	customerUUID := strings.TrimSpace(lines[9])
	if customerUUID == "" {
		return nil, fmt.Errorf("UUID клиента не может быть пустым")
	}

	request := &domain.CreateOrderTgRequest{
		Title:        title,
		Description:  description,
		WeightKg:     weightKg,
		FromCityName: fromCityName,
		FromAddress:  fromAddress,
		ToCityName:   toCityName,
		ToAddress:    toAddress,
		Price:        price,
		CustomerUUID: customerUUID,
	}

	return request, nil
}

// parseTelegramID парсит Telegram ID из строки
func parseTelegramID(text string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(text), 10, 64)
}

// parseFloat парсит float64 из строки
func parseFloat(text string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(text), 64)
}

// parseOptionalFloat парсит опциональный float64 из строки
func parseOptionalFloat(text string) *float64 {
	if text == "-" || strings.TrimSpace(text) == "" {
		return nil
	}
	if val, err := parseFloat(text); err == nil {
		return &val
	}
	return nil
}

// parseOptionalString парсит опциональную строку
func parseOptionalString(text string) *string {
	if text == "-" || strings.TrimSpace(text) == "" {
		return nil
	}
	val := strings.TrimSpace(text)
	return &val
}

// parseOptionalDate парсит опциональную дату
func parseOptionalDate(text string) *time.Time {
	if text == "-" || strings.TrimSpace(text) == "" {
		return nil
	}
	if val, err := time.Parse("2006-01-02", strings.TrimSpace(text)); err == nil {
		return &val
	}
	return nil
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
		response = "Доступные команды:\n/start - Начать работу\n/help - Показать помощь\n/status - Статус системы\n/orders - Посмотреть заказы\n/👥 Заказчики - Посмотреть заказчиков\n/🚚 Водители - Посмотреть водителей\n// Закомментировано - убираем фильтры\n// /filter - Настроить фильтры\n\nДля добавления пользователя используйте формат:\nADD_USER\nИмя\nТелефон\nTelegramID\nTelegramTag\n\nДля создания заказа используйте формат:\nADD_ORDER\nНазвание\nОписание\nВес\nОткуда город\nОткуда адрес\nКуда город\nКуда адрес\nЦена\nUUID клиента\n\nДля изменения статуса заказа используйте формат:\nARCHIVE_ORDER <UUID>\nACTIVATE_ORDER <UUID>\n\nДля настройки города и уведомлений водителя используйте формат:\nSET_CITY_AND_NOTIFICATION\nUUID, город, уведомления\n\nПримеры:\nSET_CITY_AND_NOTIFICATION\n12345678-1234-1234-1234-123456789abc, Москва, вкл\nSET_CITY_AND_NOTIFICATION\n12345678-1234-1234-1234-123456789abc, Москва, выкл\nSET_CITY_AND_NOTIFICATION\n12345678-1234-1234-1234-123456789abc, Москва\nSET_CITY_AND_NOTIFICATION\n12345678-1234-1234-1234-123456789abc, -, \nSET_CITY_AND_NOTIFICATION\n12345678-1234-1234-1234-123456789abc,, вкл\nSET_CITY_AND_NOTIFICATION\n12345678-1234-1234-1234-123456789abc,, выкл"
	case "/status":
		// Получаем статистику из базы данных
		ordersCount, err := ab.database.GetOrdersCount()
		if err != nil {
			log.Printf("Ошибка получения количества заказов: %v", err)
			ordersCount = -1
		}

		activeOrdersCount, err := ab.database.GetActiveOrdersCount()
		if err != nil {
			log.Printf("Ошибка получения количества активных заказов: %v", err)
			activeOrdersCount = -1
		}

		archivedOrdersCount := 0
		if ordersCount >= 0 && activeOrdersCount >= 0 {
			archivedOrdersCount = ordersCount - activeOrdersCount
		}

		customersCount, err := ab.database.GetCustomersCount()
		if err != nil {
			log.Printf("Ошибка получения количества клиентов: %v", err)
			customersCount = -1
		}

		driversCount, err := ab.database.GetDriversCount()
		if err != nil {
			log.Printf("Ошибка получения количества водителей: %v", err)
			driversCount = -1
		}

		if ordersCount >= 0 && customersCount >= 0 && activeOrdersCount >= 0 && driversCount >= 0 {
			response = fmt.Sprintf("✅ Система работает нормально.\n📊 Статистика:\n📋 Всего заказов: %d\n🟢 Активных: %d\n🔴 Архивных: %d\n👥 Заказчиков: %d\n🚚 Водителей: %d",
				ordersCount, activeOrdersCount, archivedOrdersCount, customersCount, driversCount)
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
	case "➕ Создать заказ":
		response = `📝 Создание нового заказа

Для создания заказа используйте следующий формат:

ADD_ORDER
Название заказа
Описание заказа
Вес (кг)
Откуда город
Откуда адрес
Куда город
Куда адрес
Цена (руб)
UUID клиента

Пример:
ADD_ORDER
Доставка мебели
Требуется доставка дивана и стола
25.5
Москва
ул. Тверская, д. 1
Санкт-Петербург
Невский проспект, д. 10
15000
12345678-1234-1234-1234-123456789abc

Отправьте сообщение с данными заказа в указанном формате.`
		keyboard = ordersMenuKeyboard()
	case "/active_orders", "🟢 Активные заказы":
		// Получаем только активные заказы
		orders, err := ab.orderService.GetActiveOrders()
		if err != nil {
			log.Printf("Ошибка получения активных заказов: %v", err)
			response = "❌ Ошибка получения активных заказов из базы данных"
		} else {
			response = ab.formatOrders(orders)
		}
		keyboard = ordersMenuKeyboard()
	case "/archived_orders", "🔴 Архивные заказы":
		// Получаем только архивные заказы
		orders, err := ab.orderService.GetOrdersByStatus("archived")
		if err != nil {
			log.Printf("Ошибка получения архивных заказов: %v", err)
			response = "❌ Ошибка получения архивных заказов из базы данных"
		} else {
			response = ab.formatOrders(orders)
		}
		keyboard = ordersMenuKeyboard()
	case "/users", "👥 Заказчики":
		// Получаем заказчиков через сервис
		customers, err := ab.customerService.GetAllCustomers()
		if err != nil {
			log.Printf("Ошибка получения заказчиков: %v", err)
			response = "❌ Ошибка получения заказчиков из базы данных"
		} else {
			response = ab.formatCustomers(customers)
		}
		keyboard = usersMenuKeyboard()
	case "/drivers", "🚚 Водители":
		// Получаем водителей через сервис
		drivers, err := ab.driverService.GetAllDrivers()
		if err != nil {
			log.Printf("Ошибка получения водителей: %v", err)
			response = "❌ Ошибка получения водителей из базы данных"
		} else {
			response = ab.formatDrivers(drivers)
		}
		keyboard = driversMenuKeyboard()

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
		keyboard = adminMainMenuKeyboard()
	default:
		// Проверяем, не является ли сообщение данными для добавления заказчика
		if strings.HasPrefix(text, "ADD_USER") {
			customer, err := ab.parseCustomerMessage(text)
			if err != nil {
				response = fmt.Sprintf("❌ Ошибка парсинга данных заказчика: %v\n\nПример правильного формата:\nADD_USER\nИван Иванов\n+79001234567\n123456789\n@ivan_username", err)
			} else {
				// Создаем заказчика через сервис
				createdCustomer, err := ab.customerService.CreateCustomer(customer.Name, customer.Phone, customer.TelegramID, customer.TelegramTag)
				if err != nil {
					response = fmt.Sprintf("❌ Ошибка создания заказчика: %v", err)
				} else {
					response = fmt.Sprintf("✅ Заказчик успешно создан!\n\n👤 Имя: %s\n📱 Телефон: %s\n🆔 Telegram ID: %s\n🏷️ Telegram Tag: %s\n🆔 UUID: %s",
						createdCustomer.Name,
						createdCustomer.Phone,
						formatTelegramID(createdCustomer.TelegramID),
						formatTelegramTag(createdCustomer.TelegramTag),
						createdCustomer.UUID)
				}
			}
		} else if strings.HasPrefix(text, "ADD_ORDER") {
			// Парсим и создаем заказ через упрощенный формат
			request, err := ab.parseOrderTgMessage(text)
			if err != nil {
				response = fmt.Sprintf("❌ Ошибка парсинга данных заказа: %v\n\nПроверьте формат и попробуйте снова.", err)
			} else {
				// Создаем заказ через сервис
				createdOrder, err := ab.orderService.CreateOrderFromTgRequest(request)
				if err != nil {
					response = fmt.Sprintf("❌ Ошибка создания заказа: %v", err)
				} else {
					response = fmt.Sprintf("✅ Заказ успешно создан!\n\n📝 %s\n📄 %s\n⚖️ %.1f кг\n🏙️ %s → %s\n💰 %.0f ₽\n🆔 ID: %s",
						createdOrder.Title,
						createdOrder.Description,
						createdOrder.WeightKg,
						*createdOrder.FromCityName,
						*createdOrder.ToCityName,
						createdOrder.Price,
						createdOrder.UUID)
				}
			}
			// Сбрасываем состояние создания заказа
			keyboard = ordersMenuKeyboard()
		} else if strings.HasPrefix(text, "SET_CITY_AND_NOTIFICATION") {
			// Парсим и выполняем команду настройки города и уведомлений водителя
			request, err := ab.parseSetCityAndNotificationMessage(text)
			if err != nil {
				response = fmt.Sprintf("❌ Ошибка парсинга команды: %v\n\nПроверьте формат и попробуйте снова.", err)
			} else {
				// Выполняем обновление через сервис
				err := ab.driverService.UpdateDriverCityAndNotifications(
					request.DriverUUID,
					request.CityName,
					request.NotificationEnabled,
				)
				if err != nil {
					response = fmt.Sprintf("❌ Ошибка обновления данных водителя: %v", err)
				} else {
					// Формируем сообщение об успехе
					var cityMsg, notificationMsg string

					if request.CityName == "-" {
						cityMsg = "город очищен"
					} else if request.CityName != "" {
						cityMsg = fmt.Sprintf("город установлен: %s", request.CityName)
					} else {
						cityMsg = "город не изменен"
					}

					if request.NotificationEnabled != nil {
						if *request.NotificationEnabled {
							notificationMsg = "уведомления включены"
						} else {
							notificationMsg = "уведомления выключены"
						}
					} else {
						notificationMsg = "уведомления не изменены"
					}

					response = fmt.Sprintf("✅ Данные водителя успешно обновлены!\n\n🚚 UUID: %s\n🏙️ %s\n🔔 %s",
						request.DriverUUID.String()[:8], cityMsg, notificationMsg)
				}
			}
		}
	}

	// Проверяем команды изменения статуса заказов
	if strings.HasPrefix(text, "ARCHIVE_ORDER ") {
		orderUUID := strings.TrimSpace(strings.TrimPrefix(text, "ARCHIVE_ORDER "))
		if orderUUID == "" {
			response = "❌ Укажите UUID заказа для архивирования\n\nПример: ARCHIVE_ORDER 12345678-1234-1234-1234-123456789abc"
		} else {
			err := ab.orderService.UpdateOrderStatus(orderUUID, "archived")
			if err != nil {
				response = fmt.Sprintf("❌ Ошибка архивирования заказа: %v", err)
			} else {
				response = fmt.Sprintf("✅ Заказ %s успешно архивирован!", orderUUID[:8])
			}
		}
	} else if strings.HasPrefix(text, "ACTIVATE_ORDER ") {
		orderUUID := strings.TrimSpace(strings.TrimPrefix(text, "ACTIVATE_ORDER "))
		if orderUUID == "" {
			response = "❌ Укажите UUID заказа для активации\n\nПример: ACTIVATE_ORDER 12345678-1234-1234-1234-123456789abc"
		} else {
			err := ab.orderService.UpdateOrderStatus(orderUUID, "active")
			if err != nil {
				response = fmt.Sprintf("❌ Ошибка активации заказа: %v", err)
			} else {
				response = fmt.Sprintf("✅ Заказ %s успешно активирован!", orderUUID[:8])
			}
		}
	}

	// splitMessage разбивает длинное сообщение на части для Telegram
	responseParts := ab.splitMessage(response, 4096) // Telegram API max message length
	for _, part := range responseParts {
		msg := tgbotapi.NewMessage(chatID, part)
		// Отключаем Markdown разметку для предотвращения ошибок
		msg.ParseMode = ""
		if keyboard.Keyboard != nil {
			msg.ReplyMarkup = keyboard
		}
		_, err := ab.bot.Send(msg)
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
	}
}

// splitMessage разбивает длинное сообщение на части для Telegram
func (ab *AdminBot) splitMessage(text string, maxLength int) []string {
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
