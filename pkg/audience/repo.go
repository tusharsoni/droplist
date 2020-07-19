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
	FindContactsByCreatedBy(ctx context.Context, createdBy string) ([]Contact, error)
	AddSegment(ctx context.Context, segment *Segment) error
	FindContactsByEmails(ctx context.Context, emails []string) ([]Contact, error)
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

func (r *sqlRepo) FindContactsByCreatedBy(ctx context.Context, createdBy string) ([]Contact, error) {
	var contacts []Contact

	err := csql.GetConn(ctx, r.db).
		Where(&Contact{CreatedBy: createdBy}).
		Find(&contacts).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query contacts", map[string]interface{}{
			"createdBy": createdBy,
		})
	}

	return contacts, nil
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
