package audience

import "time"

type List struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Name      string `gorm:"not null"`
	CreatedBy string `gorm:"not null"`
}
