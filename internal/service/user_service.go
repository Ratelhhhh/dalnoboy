package service

import (
	"fmt"
	"time"

	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"

	"github.com/google/uuid"
)

// UserService представляет сервис для работы с пользователями
type UserService struct {
	database *database.Database
}

// NewUserService создает новый экземпляр сервиса пользователей
func NewUserService(db *database.Database) *UserService {
	return &UserService{
		database: db,
	}
}

// CreateUser создает нового пользователя
func (us *UserService) CreateUser(name, phone string, telegramID *int64, telegramTag *string) (*domain.User, error) {
	// Проверяем, не существует ли уже пользователь с таким телефоном
	existingUser, err := us.database.GetUserByPhone(phone)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("пользователь с телефоном %s уже существует", phone)
	}

	// Если указан Telegram ID, проверяем, не существует ли уже пользователь с таким ID
	if telegramID != nil {
		existingUser, err := us.database.GetUserByTelegramID(*telegramID)
		if err == nil && existingUser != nil {
			return nil, fmt.Errorf("пользователь с Telegram ID %d уже существует", *telegramID)
		}
	}

	// Создаем нового пользователя
	user := &domain.User{
		UUID:        uuid.New(),
		Name:        name,
		Phone:       phone,
		TelegramID:  telegramID,
		TelegramTag: telegramTag,
		CreatedAt:   time.Now(),
	}

	// Сохраняем в базу данных
	err = us.database.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения пользователя: %v", err)
	}

	return user, nil
}

// GetAllUsers возвращает всех пользователей
func (us *UserService) GetAllUsers() ([]domain.User, error) {
	return us.database.GetAllUsers()
}

// GetUsersCount возвращает количество пользователей
func (us *UserService) GetUsersCount() (int, error) {
	return us.database.GetUsersCount()
}
