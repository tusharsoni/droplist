package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

type SaveProfileParams struct {
	AWSRegion          string `json:"aws_region" valid:"required"`
	AWSAccessKeyID     string `json:"aws_access_key_id" valid:"required"`
	AWSSecretAccessKey string `json:"aws_secret_access_key" valid:"required"`
}

type Svc interface {
	GetProfile(ctx context.Context, userUUID string) (*Profile, error)
	SaveProfile(ctx context.Context, userUUID string, p SaveProfileParams) (*Profile, error)
}

func NewSvc(repo Repo, secrets Secrets) Svc {
	return &svc{repo: repo, secrets: secrets}
}

type svc struct {
	repo    Repo
	secrets Secrets
}

func (s *svc) GetProfile(ctx context.Context, userUUID string) (*Profile, error) {
	return s.repo.GetProfileByUserUUID(ctx, userUUID)
}

func (s *svc) SaveProfile(ctx context.Context, userUUID string, p SaveProfileParams) (*Profile, error) {
	encryptedKey, err := s.secrets.Encrypt(ctx, p.AWSSecretAccessKey)
	if err != nil {
		return nil, cerror.New(err, "failed to encrypt aws secret access key", nil)
	}

	profile, err := s.GetProfile(ctx, userUUID)
	if err != nil && !cerror.HasCause(err, gorm.ErrRecordNotFound) {
		return nil, cerror.New(err, "failed to get profile", map[string]interface{}{
			"userUUID": userUUID,
		})
	}

	if profile == nil {
		profile = &Profile{
			UUID:     uuid.New().String(),
			UserUUID: userUUID,
		}
	}

	profile.AWSRegion = p.AWSRegion
	profile.AWSAccessKeyID = p.AWSAccessKeyID
	profile.AWSSecretAccessKey = encryptedKey

	err = s.repo.AddProfile(ctx, profile)
	if err != nil {
		return nil, cerror.New(err, "failed to save profile", map[string]interface{}{
			"key": encryptedKey,
		})
	}

	return profile, nil
}
