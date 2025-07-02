package db

import "database/sql/driver"

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
