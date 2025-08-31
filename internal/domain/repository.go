package domain

// OrderRepository определяет интерфейс для работы с заказами
type OrderRepository interface {
	CreateOrder(order *Order) error
	GetAllOrders() ([]Order, error)
	GetActiveOrders() ([]Order, error)
	GetOrdersByStatus(status string) ([]Order, error)
	GetOrdersByWeightRange(minWeight, maxWeight *float64) ([]Order, error)
	GetOrdersCount() (int, error)
	GetActiveOrdersCount() (int, error)
	UpdateOrderStatus(orderUUID string, status string) error
}

// CustomerRepository определяет интерфейс для работы с заказчиками
type CustomerRepository interface {
	CreateCustomer(customer *Customer) error
	GetCustomerByPhone(phone string) (*Customer, error)
	GetCustomerByTelegramID(telegramID int64) (*Customer, error)
	GetAllCustomers() ([]Customer, error)
	GetCustomersCount() (int, error)
}

// DriverRepository определяет интерфейс для работы с водителями
type DriverRepository interface {
	GetDriversCount() (int, error)
	GetAllDrivers() ([]Driver, error)
}
