package database

import (
	"database/sql"
	"fmt"
	"log"

	"dalnoboy/internal"
	"dalnoboy/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// Database представляет подключение к базе данных
type Database struct {
	DB *sql.DB
}

// Ensure Database implements OrderRepository, CustomerRepository and DriverRepository
var _ domain.OrderRepository = (*Database)(nil)
var _ domain.CustomerRepository = (*Database)(nil)
var _ domain.DriverRepository = (*Database)(nil)

// New создает новое подключение к базе данных
func New(config *internal.Config) (*Database, error) {
	connStr := config.GetDBConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия подключения к БД: %v", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	log.Println("✅ Подключение к базе данных PostgreSQL установлено успешно")

	return &Database{DB: db}, nil
}

// Close закрывает подключение к базе данных
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// GetOrdersCount возвращает количество заказов в базе данных
func (d *Database) GetOrdersCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM orders"

	err := d.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества заказов: %v", err)
	}

	return count, nil
}

// GetCustomersCount возвращает количество заказчиков в базе данных
func (d *Database) GetCustomersCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM customers"

	err := d.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества заказчиков: %v", err)
	}

	return count, nil
}

// GetAllOrders возвращает все заказы с информацией о клиентах
func (d *Database) GetAllOrders() ([]domain.Order, error) {
	query := `
		SELECT 
			o.uuid,
			o.customer_uuid,
			o.title,
			o.description,
			o.weight_kg,
			o.length_cm,
			o.width_cm,
			o.height_cm,
			o.from_city_uuid,
			o.from_address,
			o.to_city_uuid,
			o.to_address,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag,
			COALESCE(fc.name, '') as from_city_name,
			COALESCE(tc.name, '') as to_city_name
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
		LEFT JOIN cities fc ON o.from_city_uuid = fc.uuid
		LEFT JOIN cities tc ON o.to_city_uuid = tc.uuid
		ORDER BY o.created_at DESC
	`

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var tags pq.StringArray
		var fromCityName, toCityName string

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromCityUUID,
			&order.FromAddress,
			&order.ToCityUUID,
			&order.ToAddress,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
			&fromCityName,
			&toCityName,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)

		// Устанавливаем названия городов
		if fromCityName != "" {
			order.FromCityName = &fromCityName
		}
		if toCityName != "" {
			order.ToCityName = &toCityName
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return orders, nil
}

// GetActiveOrders возвращает только активные заказы
func (d *Database) GetActiveOrders() ([]domain.Order, error) {
	query := `
		SELECT 
			o.uuid,
			o.customer_uuid,
			o.title,
			o.description,
			o.weight_kg,
			o.length_cm,
			o.width_cm,
			o.height_cm,
			o.from_city_uuid,
			o.from_address,
			o.to_city_uuid,
			o.to_address,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag,
			COALESCE(fc.name, '') as from_city_name,
			COALESCE(tc.name, '') as to_city_name
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
		LEFT JOIN cities fc ON o.from_city_uuid = fc.uuid
		LEFT JOIN cities tc ON o.to_city_uuid = tc.uuid
		WHERE o.status = 'active'
		ORDER BY o.created_at DESC
	`

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var tags pq.StringArray
		var fromCityName, toCityName string

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromCityUUID,
			&order.FromAddress,
			&order.ToCityUUID,
			&order.ToAddress,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
			&fromCityName,
			&toCityName,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)

		// Устанавливаем названия городов
		if fromCityName != "" {
			order.FromCityName = &fromCityName
		}
		if toCityName != "" {
			order.ToCityName = &toCityName
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return orders, nil
}

// GetOrdersByStatus возвращает заказы по указанному статусу
func (d *Database) GetOrdersByStatus(status string) ([]domain.Order, error) {
	query := `
		SELECT 
			o.uuid,
			o.customer_uuid,
			o.title,
			o.description,
			o.weight_kg,
			o.length_cm,
			o.width_cm,
			o.height_cm,
			o.from_city_uuid,
			o.from_address,
			o.to_city_uuid,
			o.to_address,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag,
			COALESCE(fc.name, '') as from_city_name,
			COALESCE(tc.name, '') as to_city_name
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
		LEFT JOIN cities fc ON o.from_city_uuid = fc.uuid
		LEFT JOIN cities tc ON o.to_city_uuid = tc.uuid
		WHERE o.status = $1
		ORDER BY o.created_at DESC
	`

	rows, err := d.DB.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var tags pq.StringArray
		var fromCityName, toCityName string

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromCityUUID,
			&order.FromAddress,
			&order.ToCityUUID,
			&order.ToAddress,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
			&fromCityName,
			&toCityName,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)

		// Устанавливаем названия городов
		if fromCityName != "" {
			order.FromCityName = &fromCityName
		}
		if toCityName != "" {
			order.ToCityName = &toCityName
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return orders, nil
}

// GetActiveOrdersCount возвращает количество активных заказов
func (d *Database) GetActiveOrdersCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM orders WHERE status = 'active'"

	err := d.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества активных заказов: %v", err)
	}

	return count, nil
}

// UpdateOrderStatus обновляет статус заказа
func (d *Database) UpdateOrderStatus(orderUUID string, status string) error {
	query := "UPDATE orders SET status = $1 WHERE uuid = $2"

	_, err := d.DB.Exec(query, status, orderUUID)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса заказа: %v", err)
	}

	return nil
}

// GetOrdersByWeightRange возвращает заказы в указанном диапазоне веса
func (d *Database) GetOrdersByWeightRange(minWeight, maxWeight *float64) ([]domain.Order, error) {
	var query string
	var args []interface{}

	// Базовый запрос
	baseQuery := `
		SELECT 
			o.uuid,
			o.customer_uuid,
			o.title,
			o.description,
			o.weight_kg,
			o.length_cm,
			o.width_cm,
			o.height_cm,
			o.from_city_uuid,
			o.from_address,
			o.to_city_uuid,
			o.to_address,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag,
			COALESCE(fc.name, '') as from_city_name,
			COALESCE(tc.name, '') as to_city_name
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
		LEFT JOIN cities fc ON o.from_city_uuid = fc.uuid
		LEFT JOIN cities tc ON o.to_city_uuid = tc.uuid
	`

	// Формируем WHERE условие в зависимости от переданных параметров
	if minWeight != nil && maxWeight != nil {
		// Оба параметра указаны
		query = baseQuery + " WHERE o.weight_kg >= $1 AND o.weight_kg <= $2 ORDER BY o.created_at DESC"
		args = []interface{}{*minWeight, *maxWeight}
	} else if minWeight != nil {
		// Только минимальный вес
		query = baseQuery + " WHERE o.weight_kg >= $1 ORDER BY o.created_at DESC"
		args = []interface{}{*minWeight}
	} else if maxWeight != nil {
		// Только максимальный вес
		query = baseQuery + " WHERE o.weight_kg <= $1 ORDER BY o.created_at DESC"
		args = []interface{}{*maxWeight}
	} else {
		// Ни один параметр не указан - возвращаем все заказы
		query = baseQuery + " ORDER BY o.created_at DESC"
		args = []interface{}{}
	}

	rows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var tags pq.StringArray
		var fromCityName, toCityName string

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromCityUUID,
			&order.FromAddress,
			&order.ToCityUUID,
			&order.ToAddress,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
			&fromCityName,
			&toCityName,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)

		// Устанавливаем названия городов
		if fromCityName != "" {
			order.FromCityName = &fromCityName
		}
		if toCityName != "" {
			order.ToCityName = &toCityName
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return orders, nil
}

// CreateCustomer создает нового заказчика в базе данных
func (d *Database) CreateCustomer(customer *domain.Customer) error {
	query := `
		INSERT INTO customers (uuid, name, phone, telegram_id, telegram_tag, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := d.DB.Exec(query, customer.UUID.String(), customer.Name, customer.Phone, customer.TelegramID, customer.TelegramTag, customer.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка создания заказчика: %v", err)
	}

	return nil
}

// GetCustomerByPhone возвращает заказчика по номеру телефона
func (d *Database) GetCustomerByPhone(phone string) (*domain.Customer, error) {
	query := `
		SELECT uuid, name, phone, telegram_id, telegram_tag, created_at
		FROM customers
		WHERE phone = $1
	`

	var customer domain.Customer
	var uuidStr string
	err := d.DB.QueryRow(query, phone).Scan(
		&uuidStr,
		&customer.Name,
		&customer.Phone,
		&customer.TelegramID,
		&customer.TelegramTag,
		&customer.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Заказчик не найден
		}
		return nil, fmt.Errorf("ошибка получения заказчика по телефону: %v", err)
	}

	// Парсим UUID из строки
	customerUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга UUID: %v", err)
	}
	customer.UUID = customerUUID

	return &customer, nil
}

// GetCustomerByTelegramID возвращает заказчика по Telegram ID
func (d *Database) GetCustomerByTelegramID(telegramID int64) (*domain.Customer, error) {
	query := `
		SELECT uuid, name, phone, telegram_id, telegram_tag, created_at
		FROM customers
		WHERE telegram_id = $1
	`

	var customer domain.Customer
	var uuidStr string
	err := d.DB.QueryRow(query, telegramID).Scan(
		&uuidStr,
		&customer.Name,
		&customer.Phone,
		&customer.TelegramID,
		&customer.TelegramTag,
		&customer.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Заказчик не найден
		}
		return nil, fmt.Errorf("ошибка получения заказчика по Telegram ID: %v", err)
	}

	// Парсим UUID из строки
	customerUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга UUID: %v", err)
	}
	customer.UUID = customerUUID

	return &customer, nil
}

// GetAllCustomers возвращает всех заказчиков
func (d *Database) GetAllCustomers() ([]domain.Customer, error) {
	query := `
		SELECT uuid, name, phone, telegram_id, telegram_tag, created_at
		FROM customers
		ORDER BY created_at DESC
	`

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var customer domain.Customer
		var uuidStr string
		err := rows.Scan(
			&uuidStr,
			&customer.Name,
			&customer.Phone,
			&customer.TelegramID,
			&customer.TelegramTag,
			&customer.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		// Парсим UUID из строки
		customerUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга UUID: %v", err)
		}
		customer.UUID = customerUUID

		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return customers, nil
}

// CreateOrder создает новый заказ в базе данных
func (d *Database) CreateOrder(order *domain.Order) error {
	query := `
		INSERT INTO orders (
			uuid, customer_uuid, title, description, weight_kg, 
			length_cm, width_cm, height_cm, from_city_uuid, from_address, 
			to_city_uuid, to_address, tags, price, available_from, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := d.DB.Exec(query,
		order.UUID,
		order.CustomerUUID,
		order.Title,
		order.Description,
		order.WeightKg,
		order.LengthCm,
		order.WidthCm,
		order.HeightCm,
		order.FromCityUUID,
		order.FromAddress,
		order.ToCityUUID,
		order.ToAddress,
		pq.Array(order.Tags),
		order.Price,
		order.AvailableFrom,
		order.Status,
		order.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания заказа: %v", err)
	}

	return nil
}

// GetDriversCount возвращает количество водителей в базе данных
func (d *Database) GetDriversCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM drivers"

	err := d.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества водителей: %v", err)
	}

	return count, nil
}

// GetAllDrivers возвращает всех водителей с информацией о городах
func (d *Database) GetAllDrivers() ([]domain.Driver, error) {
	query := `
		SELECT 
			d.uuid,
			d.name,
			d.telegram_id,
			d.telegram_tag,
			d.notification_enabled,
			d.city_uuid,
			d.created_at,
			c.name as city_name
		FROM drivers d
		LEFT JOIN cities c ON d.city_uuid = c.uuid
		ORDER BY d.created_at DESC
	`

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var drivers []domain.Driver
	for rows.Next() {
		var driver domain.Driver
		var uuidStr string
		var cityUUIDStr sql.NullString
		err := rows.Scan(
			&uuidStr,
			&driver.Name,
			&driver.TelegramID,
			&driver.TelegramTag,
			&driver.NotificationEnabled,
			&cityUUIDStr,
			&driver.CreatedAt,
			&driver.CityName,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		// Парсим UUID из строки
		driverUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга UUID водителя: %v", err)
		}
		driver.UUID = driverUUID

		// Обрабатываем city_uuid (может быть NULL)
		if cityUUIDStr.Valid && cityUUIDStr.String != "" {
			cityUUID, err := uuid.Parse(cityUUIDStr.String)
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга UUID города: %v", err)
			}
			driver.CityUUID = &cityUUID
		} else {
			driver.CityUUID = nil
		}

		drivers = append(drivers, driver)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return drivers, nil
}

// GetCityByName возвращает город по названию
func (d *Database) GetCityByName(cityName string) (*domain.City, error) {
	query := `
		SELECT uuid, name
		FROM cities
		WHERE name = $1
	`

	var city domain.City
	var uuidStr string
	err := d.DB.QueryRow(query, cityName).Scan(&uuidStr, &city.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Город не найден
		}
		return nil, fmt.Errorf("ошибка получения города по названию: %v", err)
	}

	// Парсим UUID из строки
	cityUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга UUID города: %v", err)
	}
	city.UUID = cityUUID

	return &city, nil
}

// GetCityByUUID возвращает город по UUID
func (d *Database) GetCityByUUID(cityUUID uuid.UUID) (*domain.City, error) {
	query := `
		SELECT uuid, name
		FROM cities
		WHERE uuid = $1
	`

	var city domain.City
	var uuidStr string
	err := d.DB.QueryRow(query, cityUUID).Scan(&uuidStr, &city.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Город не найден
		}
		return nil, fmt.Errorf("ошибка получения города по UUID: %v", err)
	}

	// Парсим UUID из строки
	cityUUIDParsed, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга UUID города: %v", err)
	}
	city.UUID = cityUUIDParsed

	return &city, nil
}

// UpdateDriverCity обновляет город водителя
func (d *Database) UpdateDriverCity(driverUUID uuid.UUID, cityUUID *uuid.UUID) error {
	var query string
	var args []interface{}

	if cityUUID == nil {
		// Убираем город (устанавливаем NULL)
		query = "UPDATE drivers SET city_uuid = NULL WHERE uuid = $1"
		args = []interface{}{driverUUID}
	} else {
		// Устанавливаем город
		query = "UPDATE drivers SET city_uuid = $1 WHERE uuid = $2"
		args = []interface{}{cityUUID, driverUUID}
	}

	_, err := d.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления города водителя: %v", err)
	}

	return nil
}

// UpdateDriverNotifications обновляет статус уведомлений водителя
func (d *Database) UpdateDriverNotifications(driverUUID uuid.UUID, notificationEnabled bool) error {
	query := "UPDATE drivers SET notification_enabled = $1 WHERE uuid = $2"

	_, err := d.DB.Exec(query, notificationEnabled, driverUUID)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса уведомлений водителя: %v", err)
	}

	return nil
}

// UpdateDriverCityAndNotifications обновляет город и статус уведомлений водителя
func (d *Database) UpdateDriverCityAndNotifications(driverUUID uuid.UUID, cityUUID *uuid.UUID, notificationEnabled *bool) error {
	var query string
	var args []interface{}

	// Формируем запрос в зависимости от переданных параметров
	if cityUUID != nil && notificationEnabled != nil {
		// Обновляем и город, и уведомления
		query = "UPDATE drivers SET city_uuid = $1, notification_enabled = $2 WHERE uuid = $3"
		args = []interface{}{cityUUID, notificationEnabled, driverUUID}
	} else if cityUUID != nil {
		// Обновляем только город
		query = "UPDATE drivers SET city_uuid = $1 WHERE uuid = $2"
		args = []interface{}{cityUUID, driverUUID}
	} else if notificationEnabled != nil {
		// Обновляем только уведомления
		query = "UPDATE drivers SET notification_enabled = $1 WHERE uuid = $2"
		args = []interface{}{notificationEnabled, driverUUID}
	} else {
		// Ничего не обновляем
		return nil
	}

	_, err := d.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления данных водителя: %v", err)
	}

	return nil
}
