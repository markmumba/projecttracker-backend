package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Token     string    `gorm:"unique;not null"`
	UserID    uuid.UUID `gorm:"not null"`
	User      User      `gorm:"type:uuid;foreignKey:UserID"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
