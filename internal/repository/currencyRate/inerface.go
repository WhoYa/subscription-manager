package currencyrate

import "github.com/WhoYa/subscription-manager/pkg/db"

type CurrencyRateRepository interface {
	Create(cr *db.CurrencyRate) error
	FindByID(id string) (*db.CurrencyRate, error)
	Update(cr *db.CurrencyRate) error
	Delete(id string) error
}
