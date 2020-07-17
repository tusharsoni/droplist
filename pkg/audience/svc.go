package audience

import (
	"context"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

type Svc interface {
	CreateList(ctx context.Context, name, userUUID string) (*List, error)
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
