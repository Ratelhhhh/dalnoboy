package domain

// OrderRepository определяет интерфейс для работы с заказами
type OrderRepository interface {
	CreateOrder(order *Order) error
	GetAllOrders() ([]Order, error)
	GetActiveOrders() ([]Order, error)
	GetOrdersByStatus(status string) ([]Order, error)
	GetOrdersCount() (int, error)
	GetActiveOrdersCount() (int, error)
	UpdateOrderStatus(orderUUID string, status string) error
}

// CustomerRepository определяет интерфейс для работы с клиентами
type CustomerRepository interface {
	GetCustomersCount() (int, error)
}

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByPhone(phone string) (*User, error)
	GetUserByTelegramID(telegramID int64) (*User, error)
	GetAllUsers() ([]User, error)
	GetUsersCount() (int, error)
}
