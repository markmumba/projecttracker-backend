package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID     `gorm:"primaryKey;autoIncrement"`
	Token     string    `gorm:"unique;not null"`
	UserID    uuid.UUID      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
