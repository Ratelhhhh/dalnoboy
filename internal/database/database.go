package database

import (
	"database/sql"
	"fmt"
	"log"

	"dalnoboy/internal"
	"dalnoboy/internal/domain"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// Database представляет подключение к базе данных
type Database struct {
	DB *sql.DB
}

// Ensure Database implements OrderRepository and CustomerRepository
var _ domain.OrderRepository = (*Database)(nil)
var _ domain.CustomerRepository = (*Database)(nil)

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
			o.created_at,
			c.name as customer_name,
			c.phone as customer_phone
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
			&order.CreatedAt,
			&order.CustomerName,
			&order.CustomerPhone,
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
