package audience

import "time"

const (
	ContactStatusSubscribed = "SUBSCRIBED"
)

type Segment struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Contact struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	CreatedBy string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Status    string `gorm:"not null"`
	Params    string `gorm:"not null"`
}
