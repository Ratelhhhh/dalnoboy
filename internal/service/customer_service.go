package service

import (
	"fmt"
	"time"

	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"

	"github.com/google/uuid"
)

// CustomerService представляет сервис для работы с заказчиками
type CustomerService struct {
	database *database.Database
}

// NewCustomerService создает новый экземпляр сервиса заказчиков
func NewCustomerService(db *database.Database) *CustomerService {
	return &CustomerService{
		database: db,
	}
}

// CreateCustomer создает нового заказчика
func (cs *CustomerService) CreateCustomer(name, phone string, telegramID *int64, telegramTag *string) (*domain.Customer, error) {
	// Проверяем, не существует ли уже заказчик с таким телефоном
	existingCustomer, err := cs.database.GetCustomerByPhone(phone)
	if err == nil && existingCustomer != nil {
		return nil, fmt.Errorf("заказчик с телефоном %s уже существует", phone)
	}

	// Если указан Telegram ID, проверяем, не существует ли уже заказчик с таким ID
	if telegramID != nil {
		existingCustomer, err := cs.database.GetCustomerByTelegramID(*telegramID)
		if err == nil && existingCustomer != nil {
			return nil, fmt.Errorf("заказчик с Telegram ID %d уже существует", *telegramID)
		}
	}

	// Создаем нового заказчика
	customer := &domain.Customer{
		UUID:        uuid.New(),
		Name:        name,
		Phone:       phone,
		TelegramID:  telegramID,
		TelegramTag: telegramTag,
		CreatedAt:   time.Now(),
	}

	// Сохраняем в базу данных
	err = cs.database.CreateCustomer(customer)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения заказчика: %v", err)
	}

	return customer, nil
}

// GetAllCustomers возвращает всех заказчиков
func (cs *CustomerService) GetAllCustomers() ([]domain.Customer, error) {
	return cs.database.GetAllCustomers()
}

// GetCustomersCount возвращает количество заказчиков
func (cs *CustomerService) GetCustomersCount() (int, error) {
	return cs.database.GetCustomersCount()
}
