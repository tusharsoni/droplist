package campaign

import (
	"context"

	"github.com/tusharsoni/copper/clogger"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Stats struct {
	CampaignUUID string
	Queued       int64
	Sent         int64
	Failed       int64
	Opens        int64
	Clicks       int64
	OpenRate     float64
	ClickRate    float64
}

type Repo interface {
	AddCampaign(ctx context.Context, campaign *Campaign) error
	GetCampaignByUUID(ctx context.Context, uuid string) (*Campaign, error)
	GetStats(ctx context.Context, uuids []string) (map[string]Stats, error)
	FindCampaignsByCreatedBy(ctx context.Context, createdBy string) ([]Campaign, error)

	AddSendTask(ctx context.Context, sendTask *SendTask) error
	GetSendTaskByUUID(ctx context.Context, uuid string) (*SendTask, error)

	AddEventLog(ctx context.Context, eventLog *EventLog) error
}

func NewSQLRepo(db *gorm.DB, logger clogger.Logger) Repo {
	return &sqlRepo{
		db:     db,
		logger: logger,
	}
}

type sqlRepo struct {
	db     *gorm.DB
	logger clogger.Logger
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

func (r *sqlRepo) AddSendTask(ctx context.Context, sendTask *SendTask) error {
	err := csql.GetConn(ctx, r.db).Save(sendTask).Error
	if err != nil {
		return cerror.New(err, "failed to add sendTask", nil)
	}

	return nil
}

func (r *sqlRepo) GetSendTaskByUUID(ctx context.Context, uuid string) (*SendTask, error) {
	var sendTask SendTask

	err := csql.GetConn(ctx, r.db).
		Where(&SendTask{UUID: uuid}).
		Find(&sendTask).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query sendTask", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &sendTask, nil
}

func (r *sqlRepo) AddEventLog(ctx context.Context, eventLog *EventLog) error {
	err := csql.GetConn(ctx, r.db).Save(eventLog).Error
	if err != nil {
		return cerror.New(err, "failed to add event log", nil)
	}

	return nil
}

func (r *sqlRepo) FindCampaignsByCreatedBy(ctx context.Context, createdBy string) ([]Campaign, error) {
	var campaigns []Campaign

	err := csql.GetConn(ctx, r.db).
		Where(&Campaign{CreatedBy: createdBy}).
		Order("updated_at DESC").
		Find(&campaigns).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query campaigns", map[string]interface{}{
			"createdBy": createdBy,
		})
	}

	return campaigns, nil
}

func (r *sqlRepo) GetStats(ctx context.Context, uuids []string) (map[string]Stats, error) {
	var (
		stats = make(map[string]Stats)
		query = `
			with
			event_stats as (
				select 
					campaign_uuid,
					coalesce(sum(case when event='CLICK' then 1 else 0 end), 0) as clicks,
					coalesce(sum(case when event='OPEN' then 1 else 0 end), 0) as opens
				from campaign_event_logs
				where campaign_uuid in (?)
				group by campaign_uuid
			),
			send_stats as (
				select 
					campaign_uuid,
					coalesce(sum(case when status='QUEUED' then 1 else 0 end), 0) as queued,
					coalesce(sum(case when status='SENT' then 1 else 0 end), 0) as sent,
					coalesce(sum(case when status='FAILED' then 1 else 0 end), 0) as failed
				from campaign_send_queue 
				where campaign_uuid in (?)
				group by campaign_uuid
			)
			select 
				s.campaign_uuid,
				queued,
				sent,
				failed,
				coalesce(opens, 0) as opens,
				coalesce(clicks, 0) as clicks,
				coalesce(cast(opens as FLOAT)/cast(sent as FLOAT), 0) as open_rate,
				coalesce(cast(clicks as FLOAT)/cast(opens as FLOAT), 0) as click_rate
			from send_stats s left join event_stats e on s.campaign_uuid=e.campaign_uuid ;
		`
	)

	rows, err := csql.GetConn(ctx, r.db).
		Raw(query, uuids, uuids).
		Rows()
	if err != nil {
		return nil, cerror.New(err, "failed to query stats", map[string]interface{}{
			"uuids": uuids,
		})
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			r.logger.Error("Failed to close db rows scanner", err)
		}
	}()

	for rows.Next() {
		var statsRow Stats

		err = rows.Scan(
			&statsRow.CampaignUUID,
			&statsRow.Queued,
			&statsRow.Sent,
			&statsRow.Failed,
			&statsRow.Opens,
			&statsRow.Clicks,
			&statsRow.OpenRate,
			&statsRow.ClickRate,
		)
		if err != nil {
			return nil, cerror.New(err, "failed to scan stats row", map[string]interface{}{
				"uuids": uuids,
			})
		}

		stats[statsRow.CampaignUUID] = statsRow
	}

	return stats, nil
}
