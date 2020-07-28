package profile

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(Profile{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate profile models", nil)
	}

	return nil
}
