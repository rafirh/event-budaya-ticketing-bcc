package model

import "time"

type Fee struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	FeeType         string     `gorm:"size:50;not null;index" json:"fee_type"`
	CalculationType string     `gorm:"size:20;not null" json:"calculation_type"`
	Amount          float64    `gorm:"type:numeric(12,2);not null" json:"amount"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (Fee) TableName() string {
	return "fee_settings"
}
