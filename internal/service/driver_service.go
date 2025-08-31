package service

import (
	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

// DriverService представляет сервис для работы с водителями
type DriverService struct {
	database *database.Database
}

// NewDriverService создает новый экземпляр сервиса водителей
func NewDriverService(db *database.Database) *DriverService {
	return &DriverService{
		database: db,
	}
}

// GetAllDrivers возвращает всех водителей
func (ds *DriverService) GetAllDrivers() ([]domain.Driver, error) {
	return ds.database.GetAllDrivers()
}

// GetDriversCount возвращает количество водителей
func (ds *DriverService) GetDriversCount() (int, error) {
	return ds.database.GetDriversCount()
}

// UpdateDriverCity обновляет город водителя
func (ds *DriverService) UpdateDriverCity(driverUUID uuid.UUID, cityName string) error {
	var cityUUID *uuid.UUID

	if cityName != "" && cityName != "-" {
		// Получаем город по названию
		city, err := ds.database.GetCityByName(cityName)
		if err != nil {
			return fmt.Errorf("ошибка получения города '%s': %v", cityName, err)
		}
		if city == nil {
			return fmt.Errorf("город '%s' не найден", cityName)
		}
		cityUUID = &city.UUID
	} else if cityName == "-" {
		// Убираем город (устанавливаем NULL)
		cityUUID = nil
	}
	// Если cityName == "", то cityUUID остается nil (не изменяем)

	return ds.database.UpdateDriverCity(driverUUID, cityUUID)
}

// UpdateDriverNotifications обновляет статус уведомлений водителя
func (ds *DriverService) UpdateDriverNotifications(driverUUID uuid.UUID, notificationEnabled bool) error {
	return ds.database.UpdateDriverNotifications(driverUUID, notificationEnabled)
}

// UpdateDriverCityAndNotifications обновляет город и статус уведомлений водителя
func (ds *DriverService) UpdateDriverCityAndNotifications(driverUUID uuid.UUID, cityName string, notificationEnabled *bool) error {
	var cityUUID *uuid.UUID

	if cityName != "" && cityName != "-" {
		// Получаем город по названию
		city, err := ds.database.GetCityByName(cityName)
		if err != nil {
			return fmt.Errorf("ошибка получения города '%s': %v", cityName, err)
		}
		if city == nil {
			return fmt.Errorf("город '%s' не найден", cityName)
		}
		cityUUID = &city.UUID
	} else if cityName == "-" {
		// Убираем город (устанавливаем NULL)
		cityUUID = nil
	}
	// Если cityName == "", то cityUUID остается nil (не изменяем)

	return ds.database.UpdateDriverCityAndNotifications(driverUUID, cityUUID, notificationEnabled)
}

// GetCityByName возвращает город по названию
func (ds *DriverService) GetCityByName(cityName string) (*domain.City, error) {
	return ds.database.GetCityByName(cityName)
}
