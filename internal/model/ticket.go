package model

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	TicketCode     string     `gorm:"size:100;uniqueIndex;not null" json:"ticket_code"`
	QRCode         *string    `json:"qr_code"`
	HolderName     string     `gorm:"size:150;not null" json:"holder_name"`
	IdentityType   string     `gorm:"size:50;not null" json:"identity_type"`
	IdentityNumber string     `gorm:"size:100;not null;index" json:"identity_number"`
	HolderPhone    string     `gorm:"size:20;not null" json:"holder_phone"`
	HolderEmail    string     `gorm:"size:150;not null" json:"holder_email"`
	Notes          string     `gorm:"type:text;not null" json:"notes"`
	IsUsed         bool       `gorm:"default:false" json:"is_used"`
	UsedAt         *time.Time `json:"used_at"`
	CreatedAt      time.Time  `json:"created_at"`
	Order          Order      `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (Ticket) TableName() string {
	return "tickets"
}
