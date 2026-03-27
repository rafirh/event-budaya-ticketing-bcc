package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID                   uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PromoterID           uuid.UUID             `gorm:"type:uuid;not null;index" json:"promoter_id"`
	CategoryID           *uuid.UUID            `gorm:"type:uuid;index" json:"category_id"`
	Title                string                `gorm:"size:200;not null" json:"title"`
	Slug                 *string               `gorm:"size:200;uniqueIndex" json:"slug"`
	Summary              *string               `gorm:"size:255" json:"summary"`
	Description          *string               `json:"description"`
	Venue                *string               `gorm:"size:200" json:"venue"`
	Address              *string               `json:"address"`
	GoogleMapsURL        *string               `json:"google_maps_url"`
	Time                 *string               `gorm:"size:50" json:"time"`
	PublishedDate        *time.Time            `gorm:"type:date" json:"published_date"`
	Latitude             *float64              `gorm:"type:numeric(10,7)" json:"latitude"`
	Longitude            *float64              `gorm:"type:numeric(10,7)" json:"longitude"`
	StartDate            *time.Time            `json:"start_date"`
	EndDate              *time.Time            `json:"end_date"`
	RegistrationDeadline *time.Time            `json:"registration_deadline"`
	IsPaid               bool                  `gorm:"default:true" json:"is_paid"`
	Price                float64               `gorm:"type:numeric(12,2);default:0" json:"price"`
	Quota                int                   `gorm:"not null;default:0;check:quota >= 0" json:"quota"`
	Sold                 int                   `gorm:"default:0;check:sold >= 0" json:"sold"`
	BannerURL            *string               `json:"banner_url"`
	Status               string                `gorm:"size:20;default:draft" json:"status"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            *time.Time            `json:"updated_at"`
	Promoter             User                  `gorm:"foreignKey:PromoterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Category             *EventCategory        `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	PaymentInfo          *EventCreationPayment `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (Event) TableName() string {
	return "events"
}
