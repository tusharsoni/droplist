package audience

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(List{}, Contact{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate audience models", nil)
	}

	return nil
}
