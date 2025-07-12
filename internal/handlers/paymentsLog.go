package handlers

import (
	"errors"
	"time"

	"github.com/WhoYa/subscription-manager/internal/repository/paymentlog"
	"github.com/WhoYa/subscription-manager/internal/service"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaymentLogHandler struct {
	repo           paymentlog.PaymentLogRepository
	paymentService service.Service
}

func NewPaymentLogHandler(r paymentlog.PaymentLogRepository, paymentService service.Service) *PaymentLogHandler {
	return &PaymentLogHandler{
		repo:           r,
		paymentService: paymentService,
	}
}

func (h *PaymentLogHandler) Create(c *fiber.Ctx) error {
	userID := c.Params("userID")
	var body struct {
		SubscriptionID string  `json:"subscription_id"`
		Amount         int64   `json:"amount"` // опционально - можем рассчитать автоматически
		Currency       string  `json:"currency"`
		RateUsed       float64 `json:"rate_used"` // опционально - можем взять текущий
		PaidAt         string  `json:"paid_at"`   // ISO8601
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

	// Рассчитываем платеж чтобы получить базовую сумму и прибыль
	paymentCalc, err := h.paymentService.CalculateUserPayment(userID, body.SubscriptionID, paidAt)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "failed to calculate payment: " + err.Error()})
	}

	// Используем рассчитанные значения или переданные пользователем
	finalAmount := paymentCalc.Amount // копейки
	if body.Amount > 0 {
		finalAmount = body.Amount // если пользователь передал свою сумму
	}

	finalRate := paymentCalc.ExchangeRate
	if body.RateUsed > 0 {
		finalRate = body.RateUsed // если пользователь передал свой курс
	}

	// Рассчитываем базовую сумму и прибыль в копейках
	baseAmountKopecks := int64(paymentCalc.BaseAmount * 100)
	profitAmountKopecks := int64(paymentCalc.ProfitAmount * 100)

	pl := db.PaymentLog{
		UserID:         userID,
		SubscriptionID: body.SubscriptionID,
		Amount:         finalAmount,
		BaseAmount:     baseAmountKopecks,
		ProfitAmount:   profitAmountKopecks,
		Currency:       curr,
		RateUsed:       finalRate,
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
