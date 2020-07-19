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

type CreateSegmentParams struct {
}

type Svc interface {
	CreateContacts(ctx context.Context, userUUID string, p []CreateContactParams) ([]CreateContactResult, error)
	ListUserContacts(ctx context.Context, userUUID string) ([]Contact, error)
	GetContact(ctx context.Context, contactUUID string) (*Contact, error)
	CreateSegment(ctx context.Context, p CreateSegmentParams) (*Segment, error)
	SegmentedContacts(ctx context.Context, userUUID, segmentUUID string) ([]Contact, error)
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

func (s *svc) CreateContacts(ctx context.Context, userUUID string, createParams []CreateContactParams) ([]CreateContactResult, error) {
	results := make([]CreateContactResult, len(createParams))

	for i, p := range createParams {
		err := s.repo.AddContact(ctx, &Contact{
			UUID:      uuid.New().String(),
			CreatedBy: userUUID,
			Email:     p.Email,
			Status:    ContactStatusSubscribed,
			Params:    p.Params,
		})
		if err != nil {
			s.logger.WithTags(map[string]interface{}{
				"email":  p.Email,
				"params": p.Params,
			}).Error("Failed to save contact", err)
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

func (s *svc) SegmentedContacts(ctx context.Context, userUUID, segmentUUID string) ([]Contact, error) {
	return s.ListUserContacts(ctx, userUUID)
}

func (s *svc) ListUserContacts(ctx context.Context, userUUID string) ([]Contact, error) {
	return s.repo.FindContactsByCreatedBy(ctx, userUUID)
}

func (s *svc) CreateSegment(ctx context.Context, p CreateSegmentParams) (*Segment, error) {
	segment := &Segment{
		UUID: uuid.New().String(),
	}

	err := s.repo.AddSegment(ctx, segment)
	if err != nil {
		return nil, cerror.New(err, "failed to save segment", nil)
	}

	return segment, nil

}

func (s *svc) GetContact(ctx context.Context, contactUUID string) (*Contact, error) {
	return s.repo.GetContactByUUID(ctx, contactUUID)
}
