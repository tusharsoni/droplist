package campaign

import (
	"context"

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
}

type SvcParams struct {
	fx.In

	Repo Repo
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo: p.Repo,
	}
}

type svc struct {
	repo Repo
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
