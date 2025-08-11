package domain

import (
	"time"
)

// Order представляет доменную модель заказа
type Order struct {
	UUID          string     `json:"uuid"`
	CustomerUUID  string     `json:"customer_uuid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	WeightKg      float64    `json:"weight_kg"`
	LengthCm      *float64   `json:"length_cm"`
	WidthCm       *float64   `json:"width_cm"`
	HeightCm      *float64   `json:"height_cm"`
	FromLocation  *string    `json:"from_location"`
	ToLocation    *string    `json:"to_location"`
	Tags          []string   `json:"tags"`
	Price         float64    `json:"price"`
	AvailableFrom *time.Time `json:"available_from"`
	CreatedAt     time.Time  `json:"created_at"`
	CustomerName  string     `json:"customer_name"`
	CustomerPhone string     `json:"customer_phone"`
}
