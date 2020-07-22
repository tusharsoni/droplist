package audience

import (
	"context"
	"net/url"
	"path"

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
	ListUserContacts(ctx context.Context, userUUID string, limit, offset int) ([]Contact, error)
	GetContact(ctx context.Context, contactUUID string) (*Contact, error)
	CreateSegment(ctx context.Context, p CreateSegmentParams) (*Segment, error)
	SegmentedContacts(ctx context.Context, userUUID, segmentUUID string, limit, offset int) ([]Contact, error)
	UnsubscribeContact(ctx context.Context, uuid string) error
	UnsubscribeURL(ctx context.Context, uuid string) string
	GetContactsByEmails(ctx context.Context, emails []string) ([]Contact, error)
}

type SvcParams struct {
	fx.In

	Repo   Repo
	Config Config
	Logger clogger.Logger
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo:   p.Repo,
		config: p.Config,
		logger: p.Logger,
	}
}

type svc struct {
	repo   Repo
	config Config
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

func (s *svc) SegmentedContacts(ctx context.Context, userUUID, segmentUUID string, limit, offset int) ([]Contact, error) {
	return s.ListUserContacts(ctx, userUUID, limit, offset)
}

func (s *svc) ListUserContacts(ctx context.Context, userUUID string, limit, offset int) ([]Contact, error) {
	return s.repo.FindContactsByCreatedBy(ctx, userUUID, limit, offset)
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

func (s *svc) UnsubscribeContact(ctx context.Context, uuid string) error {
	contact, err := s.GetContact(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to get contact", map[string]interface{}{
			"uuid": uuid,
		})
	}

	contact.Status = ContactStatusUnsubscribed

	err = s.repo.AddContact(ctx, contact)
	if err != nil {
		return cerror.New(err, "failed to save contact status", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return nil
}

func (s *svc) UnsubscribeURL(ctx context.Context, uuid string) string {
	unsubURL, _ := url.Parse(path.Join("/api/audience/contacts/", uuid, "/unsubscribe"))

	return s.config.BaseURL.ResolveReference(unsubURL).String()
}

func (s *svc) GetContactsByEmails(ctx context.Context, emails []string) ([]Contact, error) {
	return s.repo.FindContactsByEmails(ctx, emails)
}
