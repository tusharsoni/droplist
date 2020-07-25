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
	UpdateTemplate(ctx context.Context, uuid string, p CreateTemplateParams) (*Template, error)
	ListUserTemplates(ctx context.Context, userUUID string) ([]Template, error)
	GeneratePreviewHTML(ctx context.Context, templateUUID string) (string, error)
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

func (s *svc) UpdateTemplate(ctx context.Context, uuid string, p CreateTemplateParams) (*Template, error) {
	tmpl, err := s.repo.GetTemplateByUUID(ctx, uuid)
	if err != nil {
		return nil, cerror.New(err, "failed to get template", map[string]interface{}{
			"uuid": uuid,
		})
	}

	tmpl.Name = p.Name
	tmpl.Subject = p.Subject
	tmpl.PreviewText = p.PreviewText
	tmpl.HTMLBody = p.HTMLBody

	err = s.repo.AddTemplate(ctx, tmpl)
	if err != nil {
		return nil, cerror.New(err, "failed to update template", nil)
	}

	return tmpl, nil
}

func (s *svc) ListUserTemplates(ctx context.Context, userUUID string) ([]Template, error) {
	return s.repo.FindTemplatesByCreatedBy(ctx, userUUID)
}

func (s *svc) GeneratePreviewHTML(ctx context.Context, templateUUID string) (string, error) {
	tmpl, err := s.GetTemplate(ctx, templateUUID)
	if err != nil {
		return "", cerror.New(err, "failed to get template", map[string]interface{}{
			"templateUUID": templateUUID,
		})
	}

	params := map[string]interface{}{
		"Contact": map[string]string{
			"FirstName": "Jane",
			"LastName":  "Doe",
		},
		"Subject":           tmpl.Subject,
		"PreviewText":       tmpl.PreviewText,
		"UnsubscribeURL":    "https://example.com/unsubscribe",
		"OpenEventImageURL": "https://dummyimage.com/10x10/ffffff/fff.png",
	}

	return generatePreviewHTML(templateUUID, tmpl.HTMLBody, params)
}
