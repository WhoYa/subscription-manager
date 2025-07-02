package migrations

import (
	"github.com/WhoYa/subscription-manager/db"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddAllTables() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250703_add_all_tables",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&db.User{},
				&db.Subscription{},
				&db.UserSubscription{},
				&db.PaymentLog{},
				&db.GlobalSettings{},
				&db.CurrencyRate{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&db.User{},
				&db.Subscription{},
				&db.UserSubscription{},
				&db.PaymentLog{},
				&db.GlobalSettings{},
				&db.CurrencyRate{},
			)
		},
	}
}
