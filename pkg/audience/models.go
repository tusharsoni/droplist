package audience

import "time"

const (
	ContactStatusSubscribed = "SUBSCRIBED"
)

type List struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Name      string `gorm:"not null"`
	CreatedBy string `gorm:"not null"`
}

type Contact struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Email  string `gorm:"not null;unique"`
	Status string `gorm:"not null"`
	Params string `gorm:"not null"`
}

type ContactListJoin struct {
	UUID        string `gorm:"primary_key"`
	ListUUID    string `gorm:"not null"`
	ContactUUID string `gorm:"not null"`
}
