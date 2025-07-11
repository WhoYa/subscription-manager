package handlers

import (
	"errors"
	"time"

	"github.com/WhoYa/subscription-manager/internal/repository/paymentlog"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaymentLogHandler struct {
	repo paymentlog.PaymentLogRepository
}

func NewPaymentLogHandler(r paymentlog.PaymentLogRepository) *PaymentLogHandler {
	return &PaymentLogHandler{repo: r}
}

func (h *PaymentLogHandler) Create(c *fiber.Ctx) error {
	userID := c.Params("userID")
	var body struct {
		SubscriptionID string  `json:"subscription_id"`
		Amount         int64   `json:"amount"`
		Currency       string  `json:"currency"`
		RateUsed       float64 `json:"rate_used"`
		PaidAt         string  `json:"paid_at"` // ISO8601
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	paidAt, err := time.Parse(time.RFC3339, body.PaidAt)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid paid_at"})
	}

	curr := db.Currency(body.Currency)

	if curr != db.USD && curr != db.EUR && curr != db.RUB {
		return c.Status(400).JSON(fiber.Map{"error": "unsupported currency"})
	}

	if body.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "amount must be > 0"})
	}

	if body.RateUsed <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "rate_used must be > 0"})
	}

	pl := db.PaymentLog{
		UserID:         userID,
		SubscriptionID: body.SubscriptionID,
		Amount:         body.Amount,
		Currency:       curr,
		RateUsed:       body.RateUsed,
		PaidAt:         paidAt,
	}

	if err := h.repo.Create(&pl); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(pl)

}

func (h *PaymentLogHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	pl, err := h.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pl)
}

func (h *PaymentLogHandler) ListByUser(c *fiber.Ctx) error {
	userID := c.Params("userID")
	fromStr, toStr := c.Query("from"), c.Query("to")
	f, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid from date"})
	}
	t, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid to date"})
	}
	logs, err := h.repo.FindByUser(userID, f, t)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(logs)
}

func (h *PaymentLogHandler) ListBySubscription(c *fiber.Ctx) error {
	subID := c.Params("subID")
	fromStr, toStr := c.Query("from"), c.Query("to")
	f, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid from date"})
	}
	t, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid to date"})
	}
	logs, err := h.repo.FindBySubscription(subID, f, t)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(logs)
}

func (h *PaymentLogHandler) ListAll(c *fiber.Ctx) error {
	fromStr, toStr := c.Query("from"), c.Query("to")
	f, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid from date"})
	}
	t, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid to date"})
	}
	logs, err := h.repo.FindAll(f, t)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(logs)
}
