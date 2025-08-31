package service

import (
	"dalnoboy/internal/database"
	"dalnoboy/internal/domain"
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
