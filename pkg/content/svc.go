package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

type CreateTemplateParams struct {
	Name        string  `json:"name" valid:"required"`
	Subject     string  `json:"subject" valid:"required"`
	PreviewText *string `json:"preview_text" valid:"optional"`
	HTMLBody    string  `json:"html_body" valid:"required"`
}

type Svc interface {
	CreateTemplate(ctx context.Context, userUUID string, p CreateTemplateParams) (*Template, error)
	GetTemplate(ctx context.Context, uuid string) (*Template, error)
	ListUserTemplates(ctx context.Context, userUUID string) ([]Template, error)
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

func (s *svc) CreateTemplate(ctx context.Context, userUUID string, p CreateTemplateParams) (*Template, error) {
	tmpl := &Template{
		UUID:        uuid.New().String(),
		CreatedBy:   userUUID,
		Name:        p.Name,
		Subject:     p.Subject,
		PreviewText: p.PreviewText,
		HTMLBody:    p.HTMLBody,
	}

	err := s.repo.AddTemplate(ctx, tmpl)
	if err != nil {
		return nil, cerror.New(err, "failed to insert template", nil)
	}

	return tmpl, nil
}

func (s *svc) GetTemplate(ctx context.Context, uuid string) (*Template, error) {
	return s.repo.GetTemplateByUUID(ctx, uuid)
}

func (s *svc) ListUserTemplates(ctx context.Context, userUUID string) ([]Template, error) {
	return s.repo.FindTemplatesByCreatedBy(ctx, userUUID)
}
