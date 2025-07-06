package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func InitialMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250703_01_initial_migration",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
                DO $$
                BEGIN
                    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'currency_enum') THEN
                        CREATE TYPE currency_enum AS ENUM ('USD','EUR','RUB');
                    END IF;
                    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ratesource_enum') THEN
                        CREATE TYPE ratesource_enum AS ENUM ('Cifra','FF');
                    END IF;
                    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pricing_mode_enum') THEN
                        CREATE TYPE pricing_mode_enum AS ENUM ('none','percent','fixed');
                    END IF;
                END$$;
            `).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`
                DROP TYPE IF EXISTS pricing_mode_enum;
                DROP TYPE IF EXISTS ratesource_enum;
                DROP TYPE IF EXISTS currency_enum;
            `).Error
		},
	}
}
