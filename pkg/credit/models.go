package credit

import "time"

type Pack struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	UserUUID        string `gorm:"not null"`
	PaymentID       string `gorm:"not null;unique"`
	PaymentComplete bool   `gorm:"not null"`

	ProductID string `gorm:"not null"`
	UseLimit  *int64
	ExpiresAt *time.Time
}

func (p Pack) TableName() string {
	return "credit_packs"
}

type UseLog struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	PackUUID     string `gorm:"not null"`
	CampaignUUID string `gorm:"not null"`
}

func (l UseLog) TableName() string {
	return "credit_use_log"
}
