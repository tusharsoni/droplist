package audience

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	AddContact(ctx context.Context, contact *Contact) error
	GetContactByUUID(ctx context.Context, uuid string) (*Contact, error)
	CountContacts(ctx context.Context, createdBy string, status *string) (int64, error)
	FindContactsByCreatedBy(ctx context.Context, createdBy string, limit, offset int) ([]Contact, error)
	AddSegment(ctx context.Context, segment *Segment) error
	FindContactsByEmails(ctx context.Context, emails []string) ([]Contact, error)
	DeleteContactsByUUIDs(ctx context.Context, uuids []string) error
	DeleteContactsByCreatedBy(ctx context.Context, createdBy string) error
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

type sqlRepo struct {
	db *gorm.DB
}

func (r *sqlRepo) AddContact(ctx context.Context, contact *Contact) error {
	err := csql.GetConn(ctx, r.db).Save(contact).Error
	if err != nil {
		return cerror.New(err, "failed to add contact", nil)
	}

	return nil
}

func (r *sqlRepo) GetContactByUUID(ctx context.Context, uuid string) (*Contact, error) {
	var contact Contact

	err := csql.GetConn(ctx, r.db).
		Where(&Contact{UUID: uuid}).
		Find(&contact).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query contact", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &contact, nil
}

func (r *sqlRepo) CountContacts(ctx context.Context, createdBy string, status *string) (int64, error) {
	var (
		count int64
		query = csql.GetConn(ctx, r.db).
			Model(&Contact{}).
			Where(&Contact{CreatedBy: createdBy})
	)

	if status != nil {
		query = query.Where(&Contact{Status: *status})
	}

	err := query.
		Count(&count).
		Error
	if err != nil {
		return 0, cerror.New(err, "failed to query contacts count", map[string]interface{}{
			"createdBy": createdBy,
			"status":    status,
		})
	}

	return count, nil
}

func (r *sqlRepo) FindContactsByCreatedBy(ctx context.Context, createdBy string, limit, offset int) ([]Contact, error) {
	var contacts []Contact

	err := csql.GetConn(ctx, r.db).
		Where(&Contact{CreatedBy: createdBy}).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&contacts).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query contacts", map[string]interface{}{
			"createdBy": createdBy,
		})
	}

	return contacts, nil
}

func (r *sqlRepo) DeleteContactsByCreatedBy(ctx context.Context, createdBy string) error {
	err := csql.GetConn(ctx, r.db).
		Where(&Contact{CreatedBy: createdBy}).
		Delete(&Contact{}).
		Error
	if err != nil {
		return cerror.New(err, "failed to delete contacts", map[string]interface{}{
			"createdBy": createdBy,
		})
	}

	return nil
}

func (r *sqlRepo) DeleteContactsByUUIDs(ctx context.Context, uuids []string) error {
	err := csql.GetConn(ctx, r.db).
		Where("uuid in (?)", uuids).
		Delete(&Contact{}).
		Error
	if err != nil {
		return cerror.New(err, "failed to delete contacts", map[string]interface{}{
			"uuids": uuids,
		})
	}

	return nil
}

func (r *sqlRepo) AddSegment(ctx context.Context, segment *Segment) error {
	err := csql.GetConn(ctx, r.db).Save(segment).Error
	if err != nil {
		return cerror.New(err, "failed to add segment", nil)
	}

	return nil
}

func (r *sqlRepo) FindContactsByEmails(ctx context.Context, emails []string) ([]Contact, error) {
	var contacts []Contact

	err := csql.GetConn(ctx, r.db).
		Where("email in (?)", emails).
		Find(&contacts).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query contacts", map[string]interface{}{
			"emails": emails,
		})
	}

	return contacts, nil
}
