package currencyrate

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type currencyRateGormRepo struct{ orm *gorm.DB }

func NewCurrencyRateRepo(db *gorm.DB) CurrencyRateRepository {
	return &currencyRateGormRepo{orm: db}
}

func (r *currencyRateGormRepo) Create(cr *db.CurrencyRate) error {
	// Генерируем UUID если он не установлен
	if cr.ID == "" {
		cr.ID = uuid.New().String()
	}

	return r.orm.Create(cr).Error
}

func (r *currencyRateGormRepo) FindByID(id string) (*db.CurrencyRate, error) {
	var cr db.CurrencyRate
	if err := r.orm.First(&cr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cr, nil
}

func (r *currencyRateGormRepo) List(limit, offset int) ([]db.CurrencyRate, error) {
	var ary []db.CurrencyRate
	err := r.orm.
		Order("fetched_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ary).Error
	return ary, err
}

func (r *currencyRateGormRepo) LatestByCurrency(currency db.Currency) (*db.CurrencyRate, error) {
	var cr db.CurrencyRate
	err := r.orm.
		Where("currency = ?", currency).
		Order("fetched_at DESC").
		First(&cr).
		Error
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

func (r *currencyRateGormRepo) Update(cr *db.CurrencyRate) error {
	return r.orm.Save(cr).Error
}

func (r *currencyRateGormRepo) Delete(id string) error {
	return r.orm.Delete(&db.CurrencyRate{}, "id = ?", id).Error
}
