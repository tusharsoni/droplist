package content

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(Template{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate content models", nil)
	}

	return nil
}
