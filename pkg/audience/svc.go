package audience

import (
	"context"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type CreateContactParams struct {
	Email  string `json:"email" valid:"required,email"`
	Params string `json:"params" valid:"required,json"`
}

type CreateContactResult struct {
	Email   string `json:"email"`
	Success bool   `json:"success"`
}

type Svc interface {
	CreateList(ctx context.Context, name, userUUID string) (*List, error)
	CreateContacts(ctx context.Context, createParams []CreateContactParams) ([]CreateContactResult, error)
}

type SvcParams struct {
	fx.In

	Repo   Repo
	Logger clogger.Logger
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo:   p.Repo,
		logger: p.Logger,
	}
}

type svc struct {
	repo   Repo
	logger clogger.Logger
}

func (s *svc) CreateList(ctx context.Context, name, userUUID string) (*List, error) {
	list := &List{
		UUID:      uuid.New().String(),
		Name:      name,
		CreatedBy: userUUID,
	}

	err := s.repo.AddList(ctx, list)
	if err != nil {
		return nil, cerror.New(err, "failed to save list", nil)
	}

	return list, nil
}

func (s *svc) CreateContacts(ctx context.Context, createParams []CreateContactParams) ([]CreateContactResult, error) {
	results := make([]CreateContactResult, len(createParams))

	for i, p := range createParams {
		err := s.repo.AddContact(ctx, &Contact{
			UUID:   uuid.New().String(),
			Email:  p.Email,
			Status: ContactStatusSubscribed,
			Params: p.Params,
		})
		if err != nil {
			s.logger.WithTags(map[string]interface{}{
				"email":  p.Email,
				"params": p.Params,
			}).Error("Failed to add contact", err)
			results[i] = CreateContactResult{
				Email:   p.Email,
				Success: false,
			}
			continue
		}

		results[i] = CreateContactResult{
			Email:   p.Email,
			Success: true,
		}
	}

	return results, nil
}
