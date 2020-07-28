package profile

import "time"

type Profile struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	UserUUID           string `gorm:"not null;unique"`
	AWSRegion          string `gorm:"not null"`
	AWSAccessKeyID     string `gorm:"not null"`
	AWSSecretAccessKey []byte `gorm:"not null"`
}
