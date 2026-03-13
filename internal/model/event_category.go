package model

import "github.com/google/uuid"

type EventCategory struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"size:100;not null" json:"name"`
}

func (EventCategory) TableName() string {
	return "event_categories"
}
