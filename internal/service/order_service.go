package service

import (
	"dalnoboy/internal/domain"
)

// OrderService определяет интерфейс для бизнес-логики заказов
type OrderService interface {
	GetAllOrders() ([]domain.Order, error)
	GetActiveOrders() ([]domain.Order, error)
	GetOrdersByStatus(status string) ([]domain.Order, error)
	GetOrdersCount() (int, error)
	GetActiveOrdersCount() (int, error)
	UpdateOrderStatus(orderUUID string, status string) error
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

// GetActiveOrders возвращает только активные заказы
func (s *orderService) GetActiveOrders() ([]domain.Order, error) {
	return s.orderRepo.GetActiveOrders()
}

// GetOrdersByStatus возвращает заказы по указанному статусу
func (s *orderService) GetOrdersByStatus(status string) ([]domain.Order, error) {
	return s.orderRepo.GetOrdersByStatus(status)
}

// GetActiveOrdersCount возвращает количество активных заказов
func (s *orderService) GetActiveOrdersCount() (int, error) {
	return s.orderRepo.GetActiveOrdersCount()
}

// UpdateOrderStatus обновляет статус заказа
func (s *orderService) UpdateOrderStatus(orderUUID string, status string) error {
	return s.orderRepo.UpdateOrderStatus(orderUUID, status)
}
