package db

import "database/sql/driver"

type RateSource string

const (
	Cifra RateSource = "Cifra" // https://cifra-bank.ru
	FF    RateSource = "FF"    // https://bankffin.kz/ru/exchange-rates
)

func (c *RateSource) Scan(value any) error {
	*c = RateSource(value.(string))
	return nil
}

func (c RateSource) Value() (driver.Value, error) {
	return string(c), nil
}
