package models

import (
	"time"

	"gorm.io/gorm"
)

type Common struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index" swaggertype:"string" format:"date-time"`
}
