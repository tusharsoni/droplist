package campaign

import "time"

const (
	StateDraft     = "DRAFT"
	StatePublished = "PUBLISHED"
)

type Campaign struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	ListUUID  string `gorm:"not null"`
	Nickname  string `gorm:"not null"`
	FromName  string `gorm:"not null"`
	FromEmail string `gorm:"not null"`
	State     string `gorm:"not null"`
}
