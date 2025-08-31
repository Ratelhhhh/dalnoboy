package app

import (
	"dalnoboy/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// WebsiteHandler обрабатывает статические файлы сайта
func (a *App) staticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Корневой путь -> index.html
	if path == "/" {
		path = "/index.html"
	}

	// Убираем начальный слеш
	filePath := strings.TrimPrefix(path, "/")

	// Проверяем что файл существует
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Определяем Content-Type по расширению
	switch filepath.Ext(filePath) {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}

	http.ServeFile(w, r, filePath)
}

// getOrdersHandler обрабатывает запросы на получение всех заказов
func (a *App) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем query параметры для фильтрации по весу
	queryParams := r.URL.Query()
	minWeightStr := queryParams.Get("min_weight")
	maxWeightStr := queryParams.Get("max_weight")

	var minWeight, maxWeight *float64
	var err error

	// Парсим минимальный вес
	if minWeightStr != "" {
		var weight float64
		weight, err = strconv.ParseFloat(minWeightStr, 64)
		if err != nil {
			http.Error(w, "Некорректный параметр min_weight", http.StatusBadRequest)
			return
		}
		minWeight = &weight
	}

	// Парсим максимальный вес
	if maxWeightStr != "" {
		var weight float64
		weight, err = strconv.ParseFloat(maxWeightStr, 64)
		if err != nil {
			http.Error(w, "Некорректный параметр max_weight", http.StatusBadRequest)
			return
		}
		maxWeight = &weight
	}

	var orders []domain.Order

	// Если указаны параметры веса, используем фильтрацию
	if minWeight != nil || maxWeight != nil {
		orders, err = a.OrderService.GetOrdersByWeightRange(minWeight, maxWeight)
	} else {
		// Иначе получаем все заказы
		orders, err = a.OrderService.GetAllOrders()
	}

	if err != nil {
		log.Printf("Ошибка получения заказов: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Преобразуем domain.Order в формат для фронтенда (как у бота)
	response := make([]map[string]interface{}, len(orders))
	for i, order := range orders {
		// Форматируем дату
		dateStr := "Не указана"
		if order.AvailableFrom != nil {
			dateStr = order.AvailableFrom.Format("02.01.2006")
		}

		// Форматируем размеры
		dimensions := "Не указаны"
		if order.LengthCm != nil && order.WidthCm != nil && order.HeightCm != nil {
			dimensions = fmt.Sprintf("%.0f×%.0f×%.0f см", *order.LengthCm, *order.WidthCm, *order.HeightCm)
		}

		// Форматируем локации для межгородских перевозок
		fromLoc := "Не указано"
		toLoc := "Не указано"

		if order.FromCityName != nil && order.ToCityName != nil {
			// Основной маршрут между городами
			fromLoc = fmt.Sprintf("%s → %s", *order.FromCityName, *order.ToCityName)

			// Адреса в одной строке
			if order.FromAddress != nil && order.ToAddress != nil {
				toLoc = fmt.Sprintf("%s: %s | %s: %s",
					*order.FromCityName, *order.FromAddress,
					*order.ToCityName, *order.ToAddress)
			} else if order.FromAddress != nil {
				toLoc = fmt.Sprintf("%s: %s", *order.FromCityName, *order.FromAddress)
			} else if order.ToAddress != nil {
				toLoc = fmt.Sprintf("%s: %s", *order.ToCityName, *order.ToAddress)
			} else {
				toLoc = "Адреса не указаны"
			}
		} else if order.FromCityName != nil {
			fromLoc = *order.FromCityName
			if order.FromAddress != nil {
				toLoc = fmt.Sprintf("Адрес: %s", *order.FromAddress)
			}
		} else if order.ToCityName != nil {
			toLoc = *order.ToCityName
			if order.ToAddress != nil {
				fromLoc = fmt.Sprintf("Адрес: %s", *order.ToAddress)
			}
		}

		// Форматируем теги
		tagsStr := "Нет тегов"
		if len(order.Tags) > 0 {
			tagsStr = strings.Join(order.Tags, ", ")
		}

		response[i] = map[string]interface{}{
			"id":          order.UUID[:8],
			"title":       order.Title,
			"description": order.Description,
			"customer":    order.CustomerName,
			"phone":       order.CustomerPhone,
			"from":        fromLoc,
			"to":          toLoc,
			"weight":      order.WeightKg,
			"dimensions":  dimensions,
			"tags":        tagsStr,
			"price":       order.Price,
			"date":        dateStr,
			"uuid":        order.UUID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
