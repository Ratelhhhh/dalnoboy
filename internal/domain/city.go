package domain

import (
	"github.com/google/uuid"
)

// City представляет доменную модель города
type City struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}
