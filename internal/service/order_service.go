package service

import (
	"fmt"
	"time"

	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"

	"github.com/google/uuid"
)

// OrderService представляет сервис для работы с заказами
type OrderService struct {
	database *database.Database
	cityRepo domain.CityRepository
}

// NewOrderService создает новый экземпляр сервиса заказов
func NewOrderService(db *database.Database, cityRepo domain.CityRepository) *OrderService {
	return &OrderService{
		database: db,
		cityRepo: cityRepo,
	}
}

// CreateOrder создает новый заказ
func (os *OrderService) CreateOrder(
	customerUUID, title, description string,
	weightKg float64,
	lengthCm, widthCm, heightCm *float64,
	fromCityUUID, fromAddress, toCityUUID, toAddress *string,
	tags []string,
	price float64,
	availableFrom *time.Time,
) (*domain.Order, error) {
	// Проверяем, что customerUUID не пустой
	if customerUUID == "" {
		return nil, fmt.Errorf("customerUUID не может быть пустым")
	}

	// Проверяем обязательные поля
	if title == "" {
		return nil, fmt.Errorf("название заказа не может быть пустым")
	}
	if description == "" {
		return nil, fmt.Errorf("описание заказа не может быть пустым")
	}
	if weightKg <= 0 {
		return nil, fmt.Errorf("вес должен быть больше нуля")
	}
	if price <= 0 {
		return nil, fmt.Errorf("цена должна быть больше нуля")
	}

	// Создаем новый заказ
	order := &domain.Order{
		UUID:          uuid.New().String(),
		CustomerUUID:  customerUUID,
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
		Status:        domain.OrderStatusActive,
		CreatedAt:     time.Now(),
	}

	// Сохраняем в базу данных
	err := os.database.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения заказа: %v", err)
	}

	return order, nil
}

// CreateOrderFromTgRequest создает новый заказ из упрощенного Telegram-запроса
func (os *OrderService) CreateOrderFromTgRequest(request *domain.CreateOrderTgRequest) (*domain.Order, error) {
	// Проверяем, что customerUUID не пустой
	if request.CustomerUUID == "" {
		return nil, fmt.Errorf("customerUUID не может быть пустым")
	}

	// Проверяем обязательные поля
	if request.Title == "" {
		return nil, fmt.Errorf("название заказа не может быть пустым")
	}
	if request.Description == "" {
		return nil, fmt.Errorf("описание заказа не может быть пустым")
	}
	if request.WeightKg <= 0 {
		return nil, fmt.Errorf("вес должен быть больше нуля")
	}
	if request.Price <= 0 {
		return nil, fmt.Errorf("цена должна быть больше нуля")
	}

	// Ищем UUID городов по названиям
	var fromCityUUID, toCityUUID *string

	if request.FromCityName != "" {
		fromCity, err := os.cityRepo.GetCityByName(request.FromCityName)
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска города отправления '%s': %v", request.FromCityName, err)
		}
		if fromCity == nil {
			return nil, fmt.Errorf("город отправления '%s' не найден в базе данных", request.FromCityName)
		}
		fromCityUUIDStr := fromCity.UUID.String()
		fromCityUUID = &fromCityUUIDStr
	}

	if request.ToCityName != "" {
		toCity, err := os.cityRepo.GetCityByName(request.ToCityName)
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска города назначения '%s': %v", request.ToCityName, err)
		}
		if toCity == nil {
			return nil, fmt.Errorf("город назначения '%s' не найден в базе данных", request.ToCityName)
		}
		toCityUUIDStr := toCity.UUID.String()
		toCityUUID = &toCityUUIDStr
	}

	// Создаем новый заказ
	order := &domain.Order{
		UUID:          uuid.New().String(),
		CustomerUUID:  request.CustomerUUID,
		Title:         request.Title,
		Description:   request.Description,
		WeightKg:      request.WeightKg,
		LengthCm:      nil, // Упрощенный формат не включает размеры
		WidthCm:       nil,
		HeightCm:      nil,
		FromCityUUID:  fromCityUUID, // Теперь используем найденные UUID
		FromAddress:   &request.FromAddress,
		FromCityName:  &request.FromCityName,
		ToCityUUID:    toCityUUID, // Теперь используем найденные UUID
		ToAddress:     &request.ToAddress,
		ToCityName:    &request.ToCityName,
		Tags:          []string{}, // Упрощенный формат не включает теги, передаем пустой массив
		Price:         request.Price,
		AvailableFrom: nil, // Упрощенный формат не включает дату
		Status:        domain.OrderStatusActive,
		CreatedAt:     time.Now(),
	}

	// Логируем создание заказа для отладки
	fmt.Printf("Создаем заказ: UUID=%s, FromCityUUID=%v, ToCityUUID=%v, Tags=%v\n",
		order.UUID, order.FromCityUUID, order.ToCityUUID, order.Tags)

	// Сохраняем в базу данных
	err := os.database.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения заказа: %v", err)
	}

	return order, nil
}

// GetAllOrders возвращает все заказы
func (os *OrderService) GetAllOrders() ([]domain.Order, error) {
	return os.database.GetAllOrders()
}

// GetOrdersCount возвращает количество заказов
func (os *OrderService) GetOrdersCount() (int, error) {
	return os.database.GetOrdersCount()
}

// GetActiveOrders возвращает только активные заказы
func (os *OrderService) GetActiveOrders() ([]domain.Order, error) {
	return os.database.GetActiveOrders()
}

// GetOrdersByStatus возвращает заказы по указанному статусу
func (os *OrderService) GetOrdersByStatus(status string) ([]domain.Order, error) {
	return os.database.GetOrdersByStatus(status)
}

// GetOrdersByWeightRange возвращает заказы в указанном диапазоне веса
func (os *OrderService) GetOrdersByWeightRange(minWeight, maxWeight *float64) ([]domain.Order, error) {
	// Валидация параметров
	if minWeight != nil && *minWeight < 0 {
		return nil, fmt.Errorf("минимальный вес не может быть отрицательным")
	}
	if maxWeight != nil && *maxWeight < 0 {
		return nil, fmt.Errorf("максимальный вес не может быть отрицательным")
	}
	if minWeight != nil && maxWeight != nil && *minWeight > *maxWeight {
		return nil, fmt.Errorf("минимальный вес не может быть больше максимального")
	}

	return os.database.GetOrdersByWeightRange(minWeight, maxWeight)
}

// UpdateOrderStatus обновляет статус заказа
func (os *OrderService) UpdateOrderStatus(orderUUID string, status string) error {
	return os.database.UpdateOrderStatus(orderUUID, status)
}
