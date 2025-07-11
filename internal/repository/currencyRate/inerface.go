package currencyrate

import "github.com/WhoYa/subscription-manager/pkg/db"

type CurrencyRateRepository interface {
	Create(cr *db.CurrencyRate) error
	FindByID(id string) (*db.CurrencyRate, error)
	List(limit, offset int) ([]db.CurrencyRate, error)
	LatestByCurrency(currency db.Currency) (*db.CurrencyRate, error)
	Update(cr *db.CurrencyRate) error
	Delete(id string) error
}
