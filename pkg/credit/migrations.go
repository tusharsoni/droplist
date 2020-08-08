package credit

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(Pack{}, UseLog{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate credit models", nil)
	}

	return nil
}
