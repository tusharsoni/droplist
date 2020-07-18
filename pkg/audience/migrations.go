package audience

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(List{}, Contact{}, ContactListJoin{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate audience models", nil)
	}

	err = db.Model(Contact{}).AddUniqueIndex("idx_created_by_email", "created_by", "email").Error
	if err != nil {
		return cerror.New(err, "failed to add unique index", map[string]interface{}{
			"table": "contacts",
			"name":  "idx_created_by_email",
		})
	}

	return nil
}
