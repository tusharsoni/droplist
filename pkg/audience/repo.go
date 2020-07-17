package audience

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	AddList(ctx context.Context, list *List) error
	AddContact(ctx context.Context, contact *Contact) error
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

type sqlRepo struct {
	db *gorm.DB
}

func (r *sqlRepo) AddList(ctx context.Context, list *List) error {
	err := csql.GetConn(ctx, r.db).Save(list).Error
	if err != nil {
		return cerror.New(err, "failed to add list", nil)
	}

	return nil
}

func (r *sqlRepo) AddContact(ctx context.Context, contact *Contact) error {
	err := csql.GetConn(ctx, r.db).Save(contact).Error
	if err != nil {
		return cerror.New(err, "failed to add contact", nil)
	}

	return nil
}
