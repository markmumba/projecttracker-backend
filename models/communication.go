package models

import (
	"time"

	"github.com/google/uuid"
)

type CommunicationHistory struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	Sender     User      `gorm:"foreignKey:SenderID"`
	Receiver   User      `gorm:"foreignKey:ReceiverID"`
}
