package database

import (
	"database/sql"
	"fmt"
	"log"

	"dalnoboy/internal"

	_ "github.com/lib/pq"
)

// Database представляет подключение к базе данных
type Database struct {
	DB *sql.DB
}

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
