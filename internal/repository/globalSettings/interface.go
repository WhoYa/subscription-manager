package globalsettings

import "github.com/WhoYa/subscription-manager/pkg/db"

type GlobalSettingsRepository interface {
	Create(gs *db.GlobalSettings) error
	Update(gs *db.GlobalSettings) error
	Get() (*db.GlobalSettings, error)
}
