package healthcheck

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
)

func Run() error {
	gormDB, err := db.Open()
	if err != nil {
		return err
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
