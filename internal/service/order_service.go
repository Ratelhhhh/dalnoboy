package service

import (
	"dalnoboy/internal/domain"
)

// OrderService определяет интерфейс для бизнес-логики заказов
type OrderService interface {
	GetAllOrders() ([]domain.Order, error)
	GetOrdersCount() (int, error)
}

// orderService реализует OrderService
type orderService struct {
	orderRepo domain.OrderRepository
}

// NewOrderService создает новый экземпляр OrderService
func NewOrderService(orderRepo domain.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

// GetAllOrders возвращает все заказы через репозиторий
func (s *orderService) GetAllOrders() ([]domain.Order, error) {
	return s.orderRepo.GetAllOrders()
}

// GetOrdersCount возвращает количество заказов через репозиторий
func (s *orderService) GetOrdersCount() (int, error) {
	return s.orderRepo.GetOrdersCount()
}
