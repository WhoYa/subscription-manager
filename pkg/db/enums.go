package db

import "database/sql/driver"

// Currency
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	RUB Currency = "RUB"
)

func (c *Currency) Scan(value any) error {
	*c = Currency(value.(string))
	return nil
}
func (c Currency) Value() (driver.Value, error) {
	return string(c), nil
}

// PricingMode
type PricingMode string

const (
	None    PricingMode = "none"
	Percent PricingMode = "percent"
	Fixed   PricingMode = "fixed"
)

func (c *PricingMode) Scan(value any) error {
	*c = PricingMode(value.(string))
	return nil
}

func (c PricingMode) Value() (driver.Value, error) {
	return string(c), nil
}

// RateSource
type RateSource string

const (
	Cifra  RateSource = "Cifra"  // https://cifra-bank.ru
	FF     RateSource = "FF"     // https://bankffin.kz/ru/exchange-rates
	Manual RateSource = "Manual" // ручной ввод админом
)

func (c *RateSource) Scan(value any) error {
	*c = RateSource(value.(string))
	return nil
}

func (c RateSource) Value() (driver.Value, error) {
	return string(c), nil
}
