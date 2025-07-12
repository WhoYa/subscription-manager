package handlers

import (
	"errors"
	"log"
	"strconv"

	repo "github.com/WhoYa/subscription-manager/internal/repository/subscription"
	dbpkg "github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	repo repo.SubscriptionRepository
}

func NewSubscriptionHandler(r repo.SubscriptionRepository) *SubscriptionHandler {
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
		log.Printf("SUBSCRIPTION: Failed to parse request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	log.Printf("SUBSCRIPTION: Creating subscription request - ServiceName: %s, BasePrice: %.2f, BaseCurrency: %s, PeriodDays: %d",
		body.ServiceName, body.BasePrice, body.BaseCurrency, body.PeriodDays)

	if exist, err := h.repo.FindByServiceName(body.ServiceName); err == nil && exist != nil {
		log.Printf("SUBSCRIPTION: Service with name '%s' already exists", body.ServiceName)
		return c.Status(409).JSON(fiber.Map{"error": "subscription with this service_name already exists"})
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("SUBSCRIPTION: Error checking existing service: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	curr := dbpkg.Currency(body.BaseCurrency)
	if curr != dbpkg.USD && curr != dbpkg.EUR {
		log.Printf("SUBSCRIPTION: Invalid currency provided: %s", body.BaseCurrency)
		return c.Status(400).JSON(fiber.Map{"error": "unsupported currency, must be USD or EUR"})
	}
	if body.PeriodDays <= 0 {
		log.Printf("SUBSCRIPTION: Invalid period days: %d", body.PeriodDays)
		return c.Status(400).JSON(fiber.Map{"error": "period_days must be > 0"})
	}

	s := dbpkg.Subscription{
		ServiceName:  body.ServiceName,
		BasePrice:    body.BasePrice,
		BaseCurrency: curr,
		PeriodDays:   body.PeriodDays,
		IsActive:     true,
	}

	log.Printf("SUBSCRIPTION: Created subscription struct: %+v", s)

	if err := h.repo.Create(&s); err != nil {
		log.Printf("SUBSCRIPTION: Failed to create subscription in database: %v", err)
		if errors.Is(err, repo.ErrDuplicateServiceName) {
			return c.Status(409).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("SUBSCRIPTION: Subscription created successfully: %+v", s)
	return c.Status(201).JSON(s)
}

func (h *SubscriptionHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid subscription id"})
	}
	s, err := h.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{"error": "subscription not found"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

func (h *SubscriptionHandler) List(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "25"))
	if err != nil || limit <= 0 {
		limit = 25
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	subs, err := h.repo.List(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(subs)
}

func (h *SubscriptionHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid subscription id"})
	}
	s, err := h.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{"error": "subscription not found"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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
		s.ServiceName = *body.ServiceName
	}
	if body.IconURL != nil {
		s.IconURL = *body.IconURL
	}
	if body.BasePrice != nil {
		s.BasePrice = *body.BasePrice
	}
	if body.BaseCurrency != nil {
		curr := dbpkg.Currency(*body.BaseCurrency)
		if curr != dbpkg.USD && curr != dbpkg.EUR {
			return c.Status(400).JSON(fiber.Map{"error": "unsupported currency"})
		}
		s.BaseCurrency = curr
	}
	if body.IsActive != nil {
		s.IsActive = *body.IsActive
	}
	if body.PeriodDays != nil {
		if *body.PeriodDays <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "period_days must be > 0"})
		}
		s.PeriodDays = *body.PeriodDays
	}

	if err := h.repo.Update(s); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

func (h *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid subscription id"})
	}
	if err := h.repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
