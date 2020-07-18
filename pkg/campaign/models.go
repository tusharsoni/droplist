package campaign

import "time"

const (
	StateDraft     = "DRAFT"
	StatePublished = "PUBLISHED"
)

const (
	SendTaskStatusQueued  = "QUEUED"
	SendTaskStatusSending = "SENDING"
	SendTaskStatusSent    = "SENT"
	SendTaskStatusFailed  = "FAILED"
)

type Campaign struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	ListUUID     string `gorm:"not null"`
	TemplateUUID string `gorm:"not null"`
	Nickname     string `gorm:"not null"`
	FromName     string `gorm:"not null"`
	FromEmail    string `gorm:"not null"`
	State        string `gorm:"not null"`
}

type SendTask struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	CampaignUUID string `gorm:"not null"`
	ContactUUID  string `gorm:"not null"`
	Status       string `gorm:"not null"`
}

func (SendTask) TableName() string {
	return "campaign_send_queue"
}
