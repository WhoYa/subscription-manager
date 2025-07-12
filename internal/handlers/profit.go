package handlers

import (
	"strconv"
	"time"

	"github.com/WhoYa/subscription-manager/internal/repository/user"
	"github.com/WhoYa/subscription-manager/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProfitHandler struct {
	profitService service.ProfitAnalytics
	userRepo      user.UserRepository
}

func NewProfitHandler(profitService service.ProfitAnalytics, userRepo user.UserRepository) *ProfitHandler {
	return &ProfitHandler{
		profitService: profitService,
		userRepo:      userRepo,
	}
}

// CheckAdminAccess проверяет является ли пользователь администратором
func (h *ProfitHandler) CheckAdminAccess(c *fiber.Ctx) error {
	userID := c.Params("adminUserID")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "admin user ID is required"})
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to check user"})
	}

	if !user.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "access denied: admin privileges required"})
	}

	return c.Next()
}

// GetMonthlyProfit возвращает прибыль за месяц
// GET /api/admin/:adminUserID/profit/monthly/:year/:month
func (h *ProfitHandler) GetMonthlyProfit(c *fiber.Ctx) error {
	yearStr := c.Params("year")
	monthStr := c.Params("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2020 || year > 2030 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid year"})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid month"})
	}

	stats, err := h.profitService.GetMonthlyProfit(year, month)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stats)
}

// GetUserProfitStats возвращает статистику прибыли по пользователям
// GET /api/admin/:adminUserID/profit/users?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z
func (h *ProfitHandler) GetUserProfitStats(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "from and to query parameters are required"})
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid from date format"})
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid to date format"})
	}

	stats, err := h.profitService.GetUserProfitStats(from, to)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stats)
}

// GetSubscriptionProfitStats возвращает статистику прибыли по подпискам
// GET /api/admin/:adminUserID/profit/subscriptions?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z
func (h *ProfitHandler) GetSubscriptionProfitStats(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "from and to query parameters are required"})
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid from date format"})
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid to date format"})
	}

	stats, err := h.profitService.GetSubscriptionProfitStats(from, to)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stats)
}

// GetTotalProfit возвращает общую прибыль за все время
// GET /api/admin/:adminUserID/profit/total
func (h *ProfitHandler) GetTotalProfit(c *fiber.Ctx) error {
	stats, err := h.profitService.GetTotalProfit()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stats)
}
