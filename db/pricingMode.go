package db

import "database/sql/driver"

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
