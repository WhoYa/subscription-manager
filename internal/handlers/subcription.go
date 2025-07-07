package handlers

import (
	"github.com/WhoYa/subscription-manager/internal/repository/subscription"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	repo subscription.SubscriptionRepository
}

func NewSubscriptionHandler(r subscription.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: r}
}

func (h *SubscriptionHandler) Create(c *fiber.Ctx) error {
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

	if err := h.repo.Create(&subscription); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(subscription)

}

func (h *SubscriptionHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	subscription, err := h.repo.FindByID(id)
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{"error": "subscription not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})

	}
	return c.JSON(subscription)
}

func (h *SubscriptionHandler) List(c *fiber.Ctx) error {
	subscriptions, err := h.repo.List(25, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(subscriptions)
}

func (h *SubscriptionHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	subscription, err := h.repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
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
	if err := h.repo.Update(subscription); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(subscription)

}
func (h *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
