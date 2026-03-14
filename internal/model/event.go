package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PromoterID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"promoter_id"`
	CategoryID    *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Slug          *string        `gorm:"size:200;uniqueIndex" json:"slug"`
	Description   *string        `json:"description"`
	Venue         *string        `gorm:"size:200" json:"venue"`
	Address       *string        `json:"address"`
	GoogleMapsURL *string        `json:"google_maps_url"`
	StartDate     *time.Time     `json:"start_date"`
	EndDate       *time.Time     `json:"end_date"`
	IsPaid        bool           `gorm:"default:true" json:"is_paid"`
	BannerURL     *string        `json:"banner_url"`
	Status        string         `gorm:"size:20;default:draft" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at"`
	Promoter      User           `gorm:"foreignKey:PromoterID" json:"-"`
	Category      *EventCategory `gorm:"foreignKey:CategoryID" json:"-"`
}

func (Event) TableName() string {
	return "events"
}
