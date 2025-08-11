package domain

// OrderRepository определяет интерфейс для работы с заказами
type OrderRepository interface {
	GetAllOrders() ([]Order, error)
	GetOrdersCount() (int, error)
}

// CustomerRepository определяет интерфейс для работы с клиентами
type CustomerRepository interface {
	GetCustomersCount() (int, error)
}
