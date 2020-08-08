package credit

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	AddPack(ctx context.Context, pack *Pack) error
	GetPackByUUID(ctx context.Context, uuid string) (*Pack, error)
	FindPacksByUserUUID(ctx context.Context, userUUID string) ([]Pack, error)
	AddUseLog(ctx context.Context, useLog *UseLog) error
	UseLogCountByPackUUIDs(ctx context.Context, packUUIDs []string) (map[string]int64, error)
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

type sqlRepo struct {
	db *gorm.DB
}

func (r *sqlRepo) AddPack(ctx context.Context, pack *Pack) error {
	err := csql.GetConn(ctx, r.db).Save(pack).Error
	if err != nil {
		return cerror.New(err, "failed to add pack", nil)
	}

	return nil
}

func (r *sqlRepo) GetPackByUUID(ctx context.Context, uuid string) (*Pack, error) {
	var pack Pack

	err := csql.GetConn(ctx, r.db).
		Where(&Pack{UUID: uuid}).
		Find(&pack).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query pack", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &pack, nil
}

func (r *sqlRepo) FindPacksByUserUUID(ctx context.Context, userUUID string) ([]Pack, error) {
	var packs []Pack

	err := csql.GetConn(ctx, r.db).
		Where(&Pack{UserUUID: userUUID}).
		Find(&packs).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query packs", map[string]interface{}{
			"userUUID": userUUID,
		})
	}

	return packs, nil
}

func (r *sqlRepo) AddUseLog(ctx context.Context, useLog *UseLog) error {
	err := csql.GetConn(ctx, r.db).Save(useLog).Error
	if err != nil {
		return cerror.New(err, "failed to add useLog", nil)
	}

	return nil
}

func (r *sqlRepo) UseLogCountByPackUUIDs(ctx context.Context, packUUIDs []string) (map[string]int64, error) {
	var (
		countByPackUUID = make(map[string]int64)
		query           = `
			SELECT pack_uuid, COUNT(*) 
				FROM credit_use_log
				WHERE pack_uuid in (?)
				GROUP BY pack_uuid
		`
	)

	rows, err := csql.GetConn(ctx, r.db).Raw(query, packUUIDs).Rows()
	if err != nil {
		return nil, cerror.New(err, "failed to run query", map[string]interface{}{
			"packUUIDs": packUUIDs,
		})
	}
	defer rows.Close()

	for rows.Next() {
		var (
			packUUID string
			count    int64
		)

		err = rows.Scan(&packUUID, &count)
		if err != nil {
			return nil, cerror.New(err, "failed to scan row", map[string]interface{}{
				"packUUIDs": packUUIDs,
			})
		}

		countByPackUUID[packUUID] = count
	}

	return countByPackUUID, nil
}
