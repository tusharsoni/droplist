package content

import "time"

type Template struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	CreatedBy   string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Subject     string `gorm:"not null"`
	PreviewText *string
	HTMLBody    string `gorm:"not null"`
}
