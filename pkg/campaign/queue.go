package campaign

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Queue interface {
	AddSendTask(ctx context.Context, sendTask *SendTask) error
	NextSendTask(ctx context.Context) (*SendTask, error)
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

func (r *sqlQueue) NextSendTask(ctx context.Context) (*SendTask, error) {
	var (
		task  SendTask
		query = `
				UPDATE campaign_send_queue SET Status='SENDING'
				WHERE uuid = (
					SELECT uuid FROM campaign_send_queue
					WHERE status='QUEUED'
					ORDER BY created_at ASC
					FOR UPDATE SKIP LOCKED 
					LIMIT 1
				)
				RETURNING *;`
	)

	err := csql.GetConn(ctx, r.db).Raw(query).Scan(&task).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, cerror.New(err, "failed to query next send task", nil)
	}

	return &task, nil
}
