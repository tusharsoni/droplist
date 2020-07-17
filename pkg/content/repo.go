package content

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	AddTemplate(ctx context.Context, template *Template) error
	GetTemplateByUUID(ctx context.Context, uuid string) (*Template, error)
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

type sqlRepo struct {
	db *gorm.DB
}

func (r *sqlRepo) AddTemplate(ctx context.Context, template *Template) error {
	err := csql.GetConn(ctx, r.db).Save(template).Error
	if err != nil {
		return cerror.New(err, "failed to add template", nil)
	}

	return nil
}

func (r *sqlRepo) GetTemplateByUUID(ctx context.Context, uuid string) (*Template, error) {
	var template Template

	err := csql.GetConn(ctx, r.db).
		Where(&Template{UUID: uuid}).
		Find(&template).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query template", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &template, nil
}
