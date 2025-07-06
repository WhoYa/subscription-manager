package handlers

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateSubscription(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ServiceName  string  `json:"service_name"`
			BasePrice    float64 `json:"base_price"`
			BaseCurrency string  `json:"base_currency"`
			PeriodDays   int     `json:"period_days"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		curr := db.Currency(body.BaseCurrency)

		if curr != db.USD && curr != db.EUR {
			return c.Status(400).JSON(fiber.Map{"error": "unsupported currency, must be USD or EUR"})
		}
		if body.PeriodDays <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "perid_days must be > 0"})
		}

		subscription := db.Subscription{
			ServiceName:  body.ServiceName,
			BasePrice:    body.BasePrice,
			BaseCurrency: curr,
			PeriodDays:   body.PeriodDays,
			IsActive:     true,
		}

		if err := dbConn.Create(&subscription).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(subscription)
	}
}

func ListSubscription(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var subscriptions []db.Subscription
		if err := dbConn.Find(&subscriptions).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(subscriptions)
	}
}

func GetSubscription(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var subscription db.Subscription
		if err := dbConn.First(&subscription, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{"error": "subscription not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(subscription)
	}
}

func UpdateSubscription(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var body struct {
			ServiceName  *string  `json:"service_name"`
			IconURL      *string  `json:"icon_url"`
			BasePrice    *float64 `json:"base_price"`
			BaseCurrency *string  `json:"base_currency"`
			IsActive     *bool    `json:"is_active"`
			PeriodDays   *int     `json:"period_days"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}
		var subscription db.Subscription
		if err := dbConn.First(&subscription, "id = ?", id).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "user not found"})
		}
		if body.ServiceName != nil {
			subscription.ServiceName = *body.ServiceName
		}
		if body.IconURL != nil {
			subscription.IconURL = *body.IconURL
		}
		if body.BasePrice != nil {
			subscription.BasePrice = *body.BasePrice
		}

		if body.BaseCurrency != nil {
			curr := db.Currency(*body.BaseCurrency)
			if curr != db.USD && curr != db.EUR {
				return c.Status(400).JSON(fiber.Map{"error": "unsupported currency"})
			}
			subscription.BaseCurrency = curr
		}
		if body.IsActive != nil {
			subscription.IsActive = *body.IsActive
		}
		if body.PeriodDays != nil {
			if *body.PeriodDays <= 0 {
				return c.Status(400).JSON(fiber.Map{"error": "period_days must be > 0"})
			}
			subscription.PeriodDays = *body.PeriodDays
		}
		if err := dbConn.Save(&subscription).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(subscription)
	}
}

func DeleteSubscription(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := dbConn.Delete(&db.Subscription{}, "id = ?", id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	}
}
