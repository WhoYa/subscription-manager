package globalsettings

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type globalSettingsGormRepo struct {
	orm *gorm.DB
}

func NewGlobalSettingsRepository(db *gorm.DB) GlobalSettingsRepository {
	return &globalSettingsGormRepo{orm: db}
}

func (r *globalSettingsGormRepo) Create(gs *db.GlobalSettings) error {
	// Генерируем UUID если он не установлен
	if gs.ID == "" {
		gs.ID = uuid.New().String()
	}

	return r.orm.Create(gs).Error
}

func (r *globalSettingsGormRepo) Update(gs *db.GlobalSettings) error {
	return r.orm.Save(gs).Error
}

func (r *globalSettingsGormRepo) Get() (*db.GlobalSettings, error) {
	var gs db.GlobalSettings
	err := r.orm.
		Order("updated_at DESC").
		First(&gs).
		Error
	if err != nil {
		return nil, err
	}
	return &gs, nil
}
