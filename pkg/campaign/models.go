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

const (
	EventOpen  = "OPEN"
	EventClick = "CLICK"
)

type Campaign struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	CreatedBy    string `gorm:"not null"`
	SegmentUUID  string `gorm:"not null"`
	TemplateUUID string `gorm:"not null"`
	Name         string `gorm:"not null"`
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
	FromName     string `gorm:"not null"`
	FromEmail    string `gorm:"not null"`
	Subject      string `gorm:"not null"`
	HTMLBody     string `gorm:"not null"`
	ToEmail      string `gorm:"not null"`
	Params       string `gorm:"not null"`
	Status       string `gorm:"not null"`
}

func (SendTask) TableName() string {
	return "campaign_send_queue"
}

type EventLog struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	CampaignUUID string `gorm:"not null"`
	ContactUUID  string `gorm:"not null"`
	Event        string `gorm:"not null"`
}

func (EventLog) TableName() string {
	return "campaign_event_logs"
}
