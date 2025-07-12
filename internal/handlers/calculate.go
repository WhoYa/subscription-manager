package handlers

import (
	"time"

	"github.com/WhoYa/subscription-manager/internal/service"
	"github.com/gofiber/fiber/v2"
)

type CalculateHandler struct {
	paymentService service.Service
}

func NewCalculateHandler(paymentService service.Service) *CalculateHandler {
	return &CalculateHandler{
		paymentService: paymentService,
	}
}

// CalculatePayment рассчитывает сумму к оплате для пользователя
// GET /api/calculate/:userID/:subscriptionID?due_date=2024-01-15
func (h *CalculateHandler) CalculatePayment(c *fiber.Ctx) error {
	userID := c.Params("userID")
	subscriptionID := c.Params("subscriptionID")

	if userID == "" || subscriptionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "userID and subscriptionID are required"})
	}

	// Парсим дату списания из query параметра
	dueDateParam := c.Query("due_date")
	var dueDate time.Time
	var err error

	if dueDateParam != "" {
		dueDate, err = time.Parse("2006-01-02", dueDateParam)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid due_date format, use YYYY-MM-DD"})
		}
	} else {
		// По умолчанию завтра
		dueDate = time.Now().AddDate(0, 0, 1)
	}

	// Рассчитываем сумму
	payment, err := h.paymentService.CalculateUserPayment(userID, subscriptionID, dueDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(payment)
}
