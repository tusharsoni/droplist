package campaign

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Queue interface {
	AddSendTask(ctx context.Context, sendTask *SendTask) error
}

func NewSQLQueue(db *gorm.DB) Queue {
	return &sqlQueue{
		db: db,
	}
}

type sqlQueue struct {
	db *gorm.DB
}

func (r *sqlQueue) AddSendTask(ctx context.Context, sendTask *SendTask) error {
	err := csql.GetConn(ctx, r.db).Save(sendTask).Error
	if err != nil {
		return cerror.New(err, "failed to add sendTask", nil)
	}

	return nil
}
