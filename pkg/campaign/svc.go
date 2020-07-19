package campaign

import (
	"context"
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
	CreateDraftCampaign(ctx context.Context, userUUID string, p CreateCampaignParams) (*Campaign, error)
	PublishCampaign(ctx context.Context, campaignUUID string) error
	CompleteSendTask(ctx context.Context, taskUUID, status string) error
}

type SvcParams struct {
	fx.In

	Repo     Repo
	Queue    Queue
	Audience audience.Svc
	Content  content.Svc
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo:     p.Repo,
		queue:    p.Queue,
		audience: p.Audience,
		content:  p.Content,
	}
}

type svc struct {
	repo     Repo
	queue    Queue
	audience audience.Svc
	content  content.Svc
}

func (s *svc) GetCampaign(ctx context.Context, campaignUUID string) (*Campaign, error) {
	return s.repo.GetCampaignByUUID(ctx, campaignUUID)
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

	contacts, err := s.audience.SegmentedContacts(ctx, campaign.CreatedBy, campaign.SegmentUUID)
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
		err = s.queue.AddSendTask(ctx, &SendTask{
			UUID:          uuid.New().String(),
			FromName:      campaign.FromName,
			FromEmail:     campaign.FromEmail,
			Subject:       tmpl.Subject,
			HTMLBody:      tmpl.HTMLBody,
			ToEmail:       contact.Email,
			ContactParams: contact.Params,
			Status:        SendTaskStatusQueued,
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
