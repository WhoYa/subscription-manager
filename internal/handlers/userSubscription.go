package handlers

import (
	"errors"
	"strconv"

	usrepo "github.com/WhoYa/subscription-manager/internal/repository/usersubscription"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
)

type UserSubscriptionHandler struct {
	repo usrepo.UserSubscriptionRepository
}

func NewUserSubscriptionHandler(r usrepo.UserSubscriptionRepository) *UserSubscriptionHandler {
	return &UserSubscriptionHandler{repo: r}
}

var validPricingModes = map[db.PricingMode]struct{}{
	db.None:    {},
	db.Percent: {},
	db.Fixed:   {},
}

func (h *UserSubscriptionHandler) Create(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var body struct {
		SubscriptionID string  `json:"subscription_id"`
		PricingMode    string  `json:"pricing_mode"`
		MarkupPercent  float64 `json:"markup_percent"`
		FixedFee       float64 `json:"fixed_fee"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	pm := db.PricingMode(body.PricingMode)
	if _, ok := validPricingModes[pm]; !ok {
		return c.Status(400).JSON(fiber.Map{"error": "pricing_mode must be one of none|percent|fixed"})
	}
	// для percent—>markup >0; для fixed—>fixed_fee >0
	switch pm {
	case db.Percent:
		if body.MarkupPercent <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "markup_percent must be > 0 for percent mode"})
		}
	case db.Fixed:
		if body.FixedFee <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "fixed_fee must be > 0 for fixed mode"})
		}
		if body.MarkupPercent != 0 {
			return c.Status(400).JSON(fiber.Map{"error": "markup_percent must be 0 for fixed mode"})
		}
	}

	us := db.UserSubscription{
		UserID:         userID,
		SubscriptionID: body.SubscriptionID,
		PricingMode:    pm,
		MarkupPercent:  body.MarkupPercent,
		FixedFee:       body.FixedFee,
	}

	if err := h.repo.Create(&us); err != nil {
		if errors.Is(err, usrepo.ErrDuplicateUserSubscription) {
			return c.Status(409).JSON(fiber.Map{"error": "user already subscribed to this service"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	full, err := h.repo.FindByID(us.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(full)
}

func (h *UserSubscriptionHandler) ListByUser(c *fiber.Ctx) error {
	userID := c.Params("userID")
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	list, err := h.repo.FindByUser(userID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

func (h *UserSubscriptionHandler) UpdateSettings(c *fiber.Ctx) error {
	id := c.Params("id")

	us, err := h.repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "subscription link not found"})
	}

	var body struct {
		PricingMode   *string  `json:"pricing_mode"`
		MarkupPercent *float64 `json:"markup_percent"`
		FixedFee      *float64 `json:"fixed_fee"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if body.PricingMode != nil {
		us.PricingMode = db.PricingMode(*body.PricingMode)
	}
	if body.MarkupPercent != nil {
		us.MarkupPercent = *body.MarkupPercent
	}
	if body.FixedFee != nil {
		us.FixedFee = *body.FixedFee
	}

	if err := h.repo.UpdateSettings(us); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(us)
}

func (h *UserSubscriptionHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
