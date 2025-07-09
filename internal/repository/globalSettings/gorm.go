package globalsettings

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"gorm.io/gorm"
)

type globalSettingsGormRepo struct {
	orm *gorm.DB
}

func NewGlobalSettings(db *gorm.DB) GlobalSettingsRepository {
	return &globalSettingsGormRepo{orm: db}
}

func (r *globalSettingsGormRepo) Create(gs *db.GlobalSettings) error {
	return r.orm.Create(gs).Error
}

func (r *globalSettingsGormRepo) Update(gs *db.GlobalSettings) error {
	return r.orm.Save(gs).Error
}
