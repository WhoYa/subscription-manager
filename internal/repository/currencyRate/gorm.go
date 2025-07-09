package currencyrate

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"gorm.io/gorm"
)

type currencyRateGormRepo struct {
	orm *gorm.DB
}

func NewCurrencyRateRepo(db *gorm.DB) CurrencyRateRepository {
	return &currencyRateGormRepo{orm: db}
}

func (r *currencyRateGormRepo) Create(cr *db.CurrencyRate) error {
	return r.orm.Create(cr).Error
}

func (r *currencyRateGormRepo) Update(cr *db.CurrencyRate) error {
	return r.orm.Save(cr).Error
}

func (r *currencyRateGormRepo) FindByID(id string) (*db.CurrencyRate, error) {
	var cr db.CurrencyRate
	err := r.orm.First(&cr, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cr, err
}

func (r currencyRateGormRepo) Delete(id string) error {
	return r.orm.Delete(&db.CurrencyRate{}, "id = ?", id).Error
}
