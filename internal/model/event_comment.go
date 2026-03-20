package model

import (
	"time"

	"github.com/google/uuid"
)

type EventComment struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"event_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id"`
	Comment   string         `gorm:"type:text;not null" json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	Event     Event          `gorm:"foreignKey:EventID" json:"-"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Parent    *EventComment  `gorm:"foreignKey:ParentID" json:"-"`
	Replies   []EventComment `gorm:"foreignKey:ParentID" json:"-"`
}

func (EventComment) TableName() string {
	return "event_comments"
}
