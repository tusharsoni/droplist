package campaign

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(Campaign{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate campaign models", nil)
	}

	return nil
}
