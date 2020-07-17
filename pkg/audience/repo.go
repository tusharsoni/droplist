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
	AddContactListJoin(ctx context.Context, contactListJoin *ContactListJoin) error
	FindContactsByListUUID(ctx context.Context, listUUID string) ([]Contact, error)
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

func (r *sqlRepo) AddContactListJoin(ctx context.Context, contactListJoin *ContactListJoin) error {
	err := csql.GetConn(ctx, r.db).Save(contactListJoin).Error
	if err != nil {
		return cerror.New(err, "failed to add contact list join", nil)
	}

	return nil
}

func (r *sqlRepo) FindContactsByListUUID(ctx context.Context, listUUID string) ([]Contact, error) {
	var contacts []Contact

	err := csql.GetConn(ctx, r.db).
		Model(&Contact{}).
		Joins("JOIN contact_list_joins ON contacts.uuid=contact_list_joins.contact_uuid").
		Where(&ContactListJoin{ListUUID: listUUID}).
		Find(&contacts).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query contacts", map[string]interface{}{
			"listUUID": listUUID,
		})
	}

	return contacts, nil
}
