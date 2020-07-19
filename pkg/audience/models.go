package audience

import (
	"encoding/json"
	"time"
)

const (
	ContactStatusSubscribed   = "SUBSCRIBED"
	ContactStatusUnsubscribed = "UNSUBSCRIBED"
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

func (c *Contact) ParamsJSON() (map[string]interface{}, error) {
	var params map[string]interface{}

	err := json.Unmarshal([]byte(c.Params), &params)
	if err != nil {
		return nil, err
	}

	return params, nil
}
