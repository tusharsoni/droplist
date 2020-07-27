package campaign

import (
	"context"
	"encoding/json"
	"net/url"
	"path"
	"shoot/pkg/audience"
	"shoot/pkg/content"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

type CreateCampaignParams struct {
	Segment      *audience.CreateSegmentParams `json:"segment" valid:"optional"`
	TemplateUUID string                        `json:"template_uuid" valid:"required,uuid"`
	Name         string                        `json:"name" valid:"required"`
	FromName     string                        `json:"from_name" valid:"required"`
	FromEmail    string                        `json:"from_email" valid:"required,email"`
}

type Svc interface {
	GetCampaign(ctx context.Context, campaignUUID string) (*Campaign, error)
	ListUserCampaigns(ctx context.Context, userUUID string) ([]Campaign, error)
	DeleteCampaign(ctx context.Context, uuid string) error
	CreateDraftCampaign(ctx context.Context, userUUID string, p CreateCampaignParams) (*Campaign, error)
	PublishCampaign(ctx context.Context, campaignUUID string) error
	CompleteSendTask(ctx context.Context, taskUUID, status string) error
	TestCampaign(ctx context.Context, campaignUUID string, recipients []string) error
	LogEvent(ctx context.Context, campaignUUID, contactUUID, event string) error
	CampaignStats(ctx context.Context, uuids []string) (map[string]Stats, error)
}

type SvcParams struct {
	fx.In

	Repo     Repo
	Queue    Queue
	Audience audience.Svc
	Content  content.Svc
	Config   Config
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo:     p.Repo,
		queue:    p.Queue,
		audience: p.Audience,
		content:  p.Content,
		config:   p.Config,
	}
}

type svc struct {
	repo     Repo
	queue    Queue
	audience audience.Svc
	content  content.Svc
	config   Config
}

func (s *svc) GetCampaign(ctx context.Context, campaignUUID string) (*Campaign, error) {
	return s.repo.GetCampaignByUUID(ctx, campaignUUID)
}

func (s *svc) ListUserCampaigns(ctx context.Context, userUUID string) ([]Campaign, error) {
	return s.repo.FindCampaignsByCreatedBy(ctx, userUUID)
}

func (s *svc) DeleteCampaign(ctx context.Context, uuid string) error {
	campaign, err := s.GetCampaign(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to get campaign", map[string]interface{}{
			"uuid": uuid,
		})
	}

	if campaign.State == StatePublished {
		return cerror.New(nil, "published campaign cannot be deleted", map[string]interface{}{
			"uuid": uuid,
		})
	}

	err = s.repo.DeleteCampaignByUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to delete campaign", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return nil
}

func (s *svc) CreateDraftCampaign(ctx context.Context, userUUID string, p CreateCampaignParams) (*Campaign, error) {
	var segmentP audience.CreateSegmentParams

	if p.Segment != nil {
		segmentP = *p.Segment
	}

	segment, err := s.audience.CreateSegment(ctx, segmentP)
	if err != nil {
		return nil, cerror.New(err, "failed to create segment", nil)
	}

	campaign := &Campaign{
		UUID:         uuid.New().String(),
		CreatedBy:    userUUID,
		SegmentUUID:  segment.UUID,
		TemplateUUID: p.TemplateUUID,
		Name:         p.Name,
		FromName:     p.FromName,
		FromEmail:    p.FromEmail,
		State:        StateDraft,
	}

	err = s.repo.AddCampaign(ctx, campaign)
	if err != nil {
		return nil, cerror.New(err, "failed to save campaign", nil)
	}

	return campaign, nil
}

func (s *svc) TestCampaign(ctx context.Context, campaignUUID string, recipients []string) error {
	campaign, err := s.GetCampaign(ctx, campaignUUID)
	if err != nil {
		return cerror.New(err, "failed to get campaign", map[string]interface{}{
			"uuid": campaignUUID,
		})
	}

	tmpl, err := s.content.GetTemplate(ctx, campaign.TemplateUUID)
	if err != nil {
		return cerror.New(err, "failed to get template", map[string]interface{}{
			"templateUUID": campaign.TemplateUUID,
		})
	}

	contacts, err := s.audience.GetContactsByEmails(ctx, recipients)
	if err != nil {
		return cerror.New(err, "failed to get contacts", map[string]interface{}{
			"emails": recipients,
		})
	}

	for _, contact := range contacts {
		contactParams, err := contact.ParamsJSON()
		if err != nil {
			return cerror.New(err, "failed to get contact params json", map[string]interface{}{
				"contactUUID": contact.UUID,
			})
		}

		params := map[string]interface{}{
			"Contact":           contactParams,
			"Subject":           tmpl.Subject,
			"UnsubscribeURL":    s.audience.UnsubscribeURL(ctx, contact.UUID),
			"OpenEventImageURL": s.GetOpenEventImageURL(campaign.UUID, contact.UUID),
		}

		paramsJ, err := json.Marshal(params)
		if err != nil {
			return cerror.New(err, "failed to marshal params as json", map[string]interface{}{
				"campaignUUID": campaign.UUID,
			})
		}

		err = s.queue.AddSendTask(ctx, &SendTask{
			UUID:         uuid.New().String(),
			CampaignUUID: campaign.UUID,
			ContactUUID:  contact.UUID,
			FromName:     campaign.FromName,
			FromEmail:    campaign.FromEmail,
			Subject:      tmpl.Subject,
			HTMLBody:     tmpl.HTMLBody,
			ToEmail:      contact.Email,
			Params:       string(paramsJ),
			Status:       SendTaskStatusQueued,
		})
		if err != nil {
			return cerror.New(err, "failed to queue send tasks", map[string]interface{}{
				"campaignUUID": campaign.UUID,
				"contactUUID":  contact.UUID,
			})
		}
	}

	return nil
}

func (s *svc) PublishCampaign(ctx context.Context, campaignUUID string) error {
	campaign, err := s.GetCampaign(ctx, campaignUUID)
	if err != nil {
		return cerror.New(err, "failed to get campaign", map[string]interface{}{
			"uuid": campaignUUID,
		})
	}

	if campaign.State != StateDraft {
		return cerror.New(nil, "campaign is not in draft state", map[string]interface{}{
			"uuid":  campaign.UUID,
			"state": campaign.State,
		})
	}

	// todo: query all contacts from this segment instead of only the first 1000
	contacts, err := s.audience.SegmentedContacts(ctx, campaign.CreatedBy, campaign.SegmentUUID, 1000, 0)
	if err != nil {
		return cerror.New(err, "failed to get segmented contacts", map[string]interface{}{
			"userUUID":    campaign.CreatedBy,
			"segmentUUID": campaign.SegmentUUID,
		})
	}

	tmpl, err := s.content.GetTemplate(ctx, campaign.TemplateUUID)
	if err != nil {
		return cerror.New(err, "failed to get template", map[string]interface{}{
			"templateUUID": campaign.TemplateUUID,
		})
	}

	for _, contact := range contacts {
		if contact.Status != audience.ContactStatusSubscribed {
			continue
		}

		contactParams, err := contact.ParamsJSON()
		if err != nil {
			return cerror.New(err, "failed to get contact params json", map[string]interface{}{
				"contactUUID": contact.UUID,
			})
		}

		params := map[string]interface{}{
			"Contact":           contactParams,
			"Subject":           tmpl.Subject,
			"UnsubscribeURL":    s.audience.UnsubscribeURL(ctx, contact.UUID),
			"OpenEventImageURL": s.GetOpenEventImageURL(campaign.UUID, contact.UUID),
		}

		paramsJ, err := json.Marshal(params)
		if err != nil {
			return cerror.New(err, "failed to marshal params as json", map[string]interface{}{
				"campaignUUID": campaign.UUID,
				"contactUUID":  contact.UUID,
			})
		}

		err = s.queue.AddSendTask(ctx, &SendTask{
			UUID:         uuid.New().String(),
			CampaignUUID: campaign.UUID,
			ContactUUID:  contact.UUID,
			FromName:     campaign.FromName,
			FromEmail:    campaign.FromEmail,
			Subject:      tmpl.Subject,
			HTMLBody:     tmpl.HTMLBody,
			ToEmail:      "success@simulator.amazonses.com",
			Params:       string(paramsJ),
			Status:       SendTaskStatusQueued,
		})
		if err != nil {
			return cerror.New(err, "failed to queue send tasks", map[string]interface{}{
				"campaignUUID": campaign.UUID,
				"contactUUID":  contact.UUID,
			})
		}
	}

	campaign.State = StatePublished

	err = s.repo.AddCampaign(ctx, campaign)
	if err != nil {
		return cerror.New(err, "failed to save campaign state", map[string]interface{}{
			"campaignUUID": campaign.UUID,
		})
	}

	return nil
}

func (s *svc) CompleteSendTask(ctx context.Context, taskUUID, status string) error {
	if status != SendTaskStatusFailed && status != SendTaskStatusSent {
		return cerror.New(nil, "task status must be failed or sent", map[string]interface{}{
			"status": status,
		})
	}

	task, err := s.repo.GetSendTaskByUUID(ctx, taskUUID)
	if err != nil {
		return cerror.New(err, "failed to get send task", map[string]interface{}{
			"taskUUID": taskUUID,
		})
	}

	task.Status = status

	err = s.repo.AddSendTask(ctx, task)
	if err != nil {
		return cerror.New(err, "failed to save task status", map[string]interface{}{
			"taskUUID": task.UUID,
			"status":   status,
		})
	}

	return nil
}

func (s *svc) LogEvent(ctx context.Context, campaignUUID, contactUUID, event string) error {
	return s.repo.AddEventLog(ctx, &EventLog{
		UUID:         uuid.New().String(),
		CampaignUUID: campaignUUID,
		ContactUUID:  contactUUID,
		Event:        event,
	})
}

func (s *svc) CampaignStats(ctx context.Context, uuids []string) (map[string]Stats, error) {
	return s.repo.GetStats(ctx, uuids)
}

func (s *svc) GetOpenEventImageURL(campaignUUID, contactUUID string) string {
	imgURL, _ := url.Parse(path.Join("/api/campaigns/", campaignUUID, "/events/", contactUUID, "/open.png"))

	return s.config.BaseURL.ResolveReference(imgURL).String()
}
