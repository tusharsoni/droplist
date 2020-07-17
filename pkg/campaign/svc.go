package campaign

import (
	"context"
	"shoot/pkg/audience"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

type CreateCampaignParams struct {
	ListUUID  string `json:"list_uuid" valid:"required,uuid"`
	Nickname  string `json:"nickname" valid:"required"`
	FromName  string `json:"from_name" valid:"required"`
	FromEmail string `json:"from_email" valid:"required,email"`
}

type Svc interface {
	CreateDraftCampaign(ctx context.Context, p CreateCampaignParams) (*Campaign, error)
	PublishCampaign(ctx context.Context, campaignUUID string) error
}

type SvcParams struct {
	fx.In

	Repo     Repo
	Queue    Queue
	Audience audience.Svc
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo:     p.Repo,
		queue:    p.Queue,
		audience: p.Audience,
	}
}

type svc struct {
	repo     Repo
	queue    Queue
	audience audience.Svc
}

func (s *svc) GetCampaign(ctx context.Context, campaignUUID string) (*Campaign, error) {
	return s.repo.GetCampaignByUUID(ctx, campaignUUID)
}

func (s *svc) CreateDraftCampaign(ctx context.Context, p CreateCampaignParams) (*Campaign, error) {
	campaign := &Campaign{
		UUID:      uuid.New().String(),
		ListUUID:  p.ListUUID,
		Nickname:  p.Nickname,
		FromName:  p.FromName,
		FromEmail: p.FromEmail,
		State:     StateDraft,
	}

	err := s.repo.AddCampaign(ctx, campaign)
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

	contacts, err := s.audience.ListContacts(ctx, campaign.ListUUID)
	if err != nil {
		return cerror.New(err, "failed to get list contacts", map[string]interface{}{
			"listUUID": campaign.ListUUID,
		})
	}

	for _, contact := range contacts {
		err = s.queue.AddSendTask(ctx, &SendTask{
			UUID:         uuid.New().String(),
			CampaignUUID: campaign.UUID,
			ContactUUID:  contact.UUID,
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
