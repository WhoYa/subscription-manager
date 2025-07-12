package handlers

import (
	"errors"
	"net/http"
	"time"

	crRepo "github.com/WhoYa/subscription-manager/internal/repository/currencyrate"
	userRepo "github.com/WhoYa/subscription-manager/internal/repository/user"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminHandler struct {
	userRepo     userRepo.UserRepository
	currencyRepo crRepo.CurrencyRateRepository
}

func NewAdminHandler(uRepo userRepo.UserRepository, crRepo crRepo.CurrencyRateRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:     uRepo,
		currencyRepo: crRepo,
	}
}

// CheckAdminAccess middleware для проверки админских прав
func (h *AdminHandler) CheckAdminAccess(c *fiber.Ctx) error {
	adminUserID := c.Params("adminUserID")
	if adminUserID == "" {
		return fiber.NewError(http.StatusBadRequest, "admin user ID required")
	}

	user, err := h.userRepo.FindByID(adminUserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(http.StatusNotFound, "admin user not found")
	} else if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if !user.IsAdmin {
		return fiber.NewError(http.StatusForbidden, "admin access required")
	}

	return c.Next()
}

// SetManualRate быстрый ввод курса валюты админом
func (h *AdminHandler) SetManualRate(c *fiber.Ctx) error {
	var body struct {
		Currency string  `json:"currency"`
		Rate     float64 `json:"rate"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid JSON")
	}

	// Валидация валюты
	curr := db.Currency(body.Currency)
	switch curr {
	case db.USD, db.EUR, db.RUB:
	default:
		return fiber.NewError(http.StatusBadRequest, "supported currencies: USD, EUR, RUB")
	}

	// Валидация курса
	if body.Rate <= 0 {
		return fiber.NewError(http.StatusBadRequest, "rate must be greater than 0")
	}

	// Создаем новый курс
	now := time.Now()
	currencyRate := &db.CurrencyRate{
		Currency:  curr,
		Value:     body.Rate,
		Source:    db.Manual,
		FetchedAt: now,
	}

	if err := h.currencyRepo.Create(currencyRate); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":  "Manual currency rate set successfully",
		"currency": body.Currency,
		"rate":     body.Rate,
		"source":   "Manual",
		"set_at":   now,
		"id":       currencyRate.ID,
	})
}

// GetCurrentRates получение текущих курсов всех валют
func (h *AdminHandler) GetCurrentRates(c *fiber.Ctx) error {
	currencies := []db.Currency{db.USD, db.EUR, db.RUB}
	rates := make(map[string]interface{})

	for _, currency := range currencies {
		rate, err := h.currencyRepo.LatestByCurrency(currency)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rates[string(currency)] = map[string]interface{}{
				"available": false,
				"message":   "No rate available",
			}
		} else if err != nil {
			rates[string(currency)] = map[string]interface{}{
				"available": false,
				"error":     err.Error(),
			}
		} else {
			rates[string(currency)] = map[string]interface{}{
				"available": true,
				"rate":      rate.Value,
				"source":    string(rate.Source),
				"set_at":    rate.FetchedAt,
				"age_hours": time.Since(rate.FetchedAt).Hours(),
			}
		}
	}

	return c.JSON(fiber.Map{
		"rates":      rates,
		"checked_at": time.Now(),
	})
}

// SetMultipleRates установка курсов для нескольких валют одновременно
func (h *AdminHandler) SetMultipleRates(c *fiber.Ctx) error {
	var body struct {
		Rates []struct {
			Currency string  `json:"currency"`
			Rate     float64 `json:"rate"`
		} `json:"rates"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid JSON")
	}

	if len(body.Rates) == 0 {
		return fiber.NewError(http.StatusBadRequest, "at least one rate required")
	}

	now := time.Now()
	results := make([]map[string]interface{}, 0, len(body.Rates))

	for _, rateData := range body.Rates {
		// Валидация валюты
		curr := db.Currency(rateData.Currency)
		switch curr {
		case db.USD, db.EUR, db.RUB:
		default:
			results = append(results, map[string]interface{}{
				"currency": rateData.Currency,
				"success":  false,
				"error":    "unsupported currency",
			})
			continue
		}

		// Валидация курса
		if rateData.Rate <= 0 {
			results = append(results, map[string]interface{}{
				"currency": rateData.Currency,
				"success":  false,
				"error":    "rate must be greater than 0",
			})
			continue
		}

		// Создаем курс
		currencyRate := &db.CurrencyRate{
			Currency:  curr,
			Value:     rateData.Rate,
			Source:    db.Manual,
			FetchedAt: now,
		}

		if err := h.currencyRepo.Create(currencyRate); err != nil {
			results = append(results, map[string]interface{}{
				"currency": rateData.Currency,
				"success":  false,
				"error":    err.Error(),
			})
		} else {
			results = append(results, map[string]interface{}{
				"currency": rateData.Currency,
				"success":  true,
				"rate":     rateData.Rate,
				"id":       currencyRate.ID,
			})
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":   "Bulk rate setting completed",
		"results":   results,
		"processed": len(body.Rates),
		"set_at":    now,
	})
}
