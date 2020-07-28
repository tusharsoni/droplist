package profile

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	GetProfileByUserUUID(ctx context.Context, userUUID string) (*Profile, error)
	AddProfile(ctx context.Context, profile *Profile) error
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

type sqlRepo struct {
	db *gorm.DB
}

func (r *sqlRepo) GetProfileByUserUUID(ctx context.Context, userUUID string) (*Profile, error) {
	var profile Profile

	err := csql.GetConn(ctx, r.db).
		Where(&Profile{UserUUID: userUUID}).
		Find(&profile).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query profile", map[string]interface{}{
			"userUUID": userUUID,
		})
	}

	return &profile, nil
}

func (r *sqlRepo) AddProfile(ctx context.Context, profile *Profile) error {
	err := csql.GetConn(ctx, r.db).Save(profile).Error
	if err != nil {
		return cerror.New(err, "failed to add profile", nil)
	}

	return nil
}
