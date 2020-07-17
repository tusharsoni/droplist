package campaign

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	AddCampaign(ctx context.Context, campaign *Campaign) error
	GetCampaignByUUID(ctx context.Context, uuid string) (*Campaign, error)
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

type sqlRepo struct {
	db *gorm.DB
}

func (r *sqlRepo) AddCampaign(ctx context.Context, campaign *Campaign) error {
	err := csql.GetConn(ctx, r.db).Save(campaign).Error
	if err != nil {
		return cerror.New(err, "failed to add campaign", nil)
	}

	return nil
}

func (r *sqlRepo) GetCampaignByUUID(ctx context.Context, uuid string) (*Campaign, error) {
	var campaign Campaign

	err := csql.GetConn(ctx, r.db).
		Where(&Campaign{UUID: uuid}).
		Find(&campaign).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query campaign", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &campaign, nil
}
