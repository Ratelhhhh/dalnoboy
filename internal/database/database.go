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

// Ensure Database implements OrderRepository, CustomerRepository and UserRepository
var _ domain.OrderRepository = (*Database)(nil)
var _ domain.CustomerRepository = (*Database)(nil)
var _ domain.UserRepository = (*Database)(nil)

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

// GetCustomersCount возвращает количество клиентов в базе данных
func (d *Database) GetCustomersCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM customers"

	err := d.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества клиентов: %v", err)
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
			o.from_location,
			o.to_location,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
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

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromLocation,
			&order.ToLocation,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)
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
			o.from_location,
			o.to_location,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
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

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromLocation,
			&order.ToLocation,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)
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
			o.from_location,
			o.to_location,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
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

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromLocation,
			&order.ToLocation,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)
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
			o.from_location,
			o.to_location,
			o.tags,
			o.price,
			o.available_from,
			o.status,
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone,
			c.telegram_id as customer_telegram_id,
			c.telegram_tag as customer_telegram_tag
		FROM orders o
		JOIN customers c ON o.customer_uuid = c.uuid
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

		err := rows.Scan(
			&order.UUID,
			&order.CustomerUUID,
			&order.Title,
			&order.Description,
			&order.WeightKg,
			&order.LengthCm,
			&order.WidthCm,
			&order.HeightCm,
			&order.FromLocation,
			&order.ToLocation,
			&tags,
			&order.Price,
			&order.AvailableFrom,
			&order.Status,
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerTelegramID,
			&order.CustomerTelegramTag,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		order.Tags = []string(tags)
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return orders, nil
}

// CreateUser создает нового пользователя в базе данных
func (d *Database) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO customers (uuid, name, phone, telegram_id, telegram_tag, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := d.DB.Exec(query, user.UUID.String(), user.Name, user.Phone, user.TelegramID, user.TelegramTag, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %v", err)
	}

	return nil
}

// GetUserByPhone возвращает пользователя по номеру телефона
func (d *Database) GetUserByPhone(phone string) (*domain.User, error) {
	query := `
		SELECT uuid, name, phone, telegram_id, telegram_tag, created_at
		FROM customers
		WHERE phone = $1
	`

	var user domain.User
	var uuidStr string
	err := d.DB.QueryRow(query, phone).Scan(
		&uuidStr,
		&user.Name,
		&user.Phone,
		&user.TelegramID,
		&user.TelegramTag,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("ошибка получения пользователя по телефону: %v", err)
	}

	// Парсим UUID из строки
	userUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга UUID: %v", err)
	}
	user.UUID = userUUID

	return &user, nil
}

// GetUserByTelegramID возвращает пользователя по Telegram ID
func (d *Database) GetUserByTelegramID(telegramID int64) (*domain.User, error) {
	query := `
		SELECT uuid, name, phone, telegram_id, telegram_tag, created_at
		FROM customers
		WHERE telegram_id = $1
	`

	var user domain.User
	var uuidStr string
	err := d.DB.QueryRow(query, telegramID).Scan(
		&uuidStr,
		&user.Name,
		&user.Phone,
		&user.TelegramID,
		&user.TelegramTag,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("ошибка получения пользователя по Telegram ID: %v", err)
	}

	// Парсим UUID из строки
	userUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга UUID: %v", err)
	}
	user.UUID = userUUID

	return &user, nil
}

// GetAllUsers возвращает всех пользователей
func (d *Database) GetAllUsers() ([]domain.User, error) {
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

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var uuidStr string
		err := rows.Scan(
			&uuidStr,
			&user.Name,
			&user.Phone,
			&user.TelegramID,
			&user.TelegramTag,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		// Парсим UUID из строки
		userUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга UUID: %v", err)
		}
		user.UUID = userUUID

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return users, nil
}

// GetUsersCount возвращает количество пользователей в базе данных
func (d *Database) GetUsersCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM customers"

	err := d.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества пользователей: %v", err)
	}

	return count, nil
}

// CreateOrder создает новый заказ в базе данных
func (d *Database) CreateOrder(order *domain.Order) error {
	query := `
		INSERT INTO orders (
			uuid, customer_uuid, title, description, weight_kg, 
			length_cm, width_cm, height_cm, from_location, to_location, 
			tags, price, available_from, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
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
		order.FromLocation,
		order.ToLocation,
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
