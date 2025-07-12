package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	repo "github.com/WhoYa/subscription-manager/internal/repository/currencyrate"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CurrencyRateHandler struct {
	repo repo.CurrencyRateRepository
}

func NewCurrencyRateHandler(r repo.CurrencyRateRepository) *CurrencyRateHandler {
	return &CurrencyRateHandler{repo: r}
}

func (h *CurrencyRateHandler) Create(c *fiber.Ctx) error {
	var body struct {
		Currency  string  `json:"currency"`
		Value     float64 `json:"value"`
		Source    string  `json:"source"`
		FetchedAt string  `json:"fetched_at"` // optional ISO8601
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid JSON")
	}

	curr := db.Currency(body.Currency)
	switch curr {
	case db.USD, db.EUR, db.RUB:
	default:
		return fiber.NewError(http.StatusBadRequest, "unsupported currency")
	}
	if body.Value <= 0 {
		return fiber.NewError(http.StatusBadRequest, "value must be > 0")
	}
	src := db.RateSource(body.Source)
	switch src {
	case db.Cifra, db.FF, db.Manual:
	default:
		return fiber.NewError(http.StatusBadRequest, "unsupported source")
	}

	var fetched time.Time
	if body.FetchedAt != "" {
		t, err := time.Parse(time.RFC3339, body.FetchedAt)
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, "invalid fetched_at")
		}
		fetched = t
	} else {
		fetched = time.Now().UTC()
	}

	cr := db.CurrencyRate{
		Currency:  curr,
		Value:     body.Value,
		Source:    src,
		FetchedAt: fetched,
	}
	if err := h.repo.Create(&cr); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusCreated).JSON(cr)
}

// Get by ID
func (h *CurrencyRateHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	cr, err := h.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(http.StatusNotFound, "currency rate not found")
	} else if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(cr)
}

// List with pagination
func (h *CurrencyRateHandler) List(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "25"))
	if err != nil || limit <= 0 {
		limit = 25
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	ary, err := h.repo.List(limit, offset)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(ary)
}

// Latest by currency
func (h *CurrencyRateHandler) Latest(c *fiber.Ctx) error {
	curr := db.Currency(c.Params("currency"))
	switch curr {
	case db.USD, db.EUR, db.RUB:
	default:
		return fiber.NewError(http.StatusBadRequest, "unsupported currency")
	}
	cr, err := h.repo.LatestByCurrency(curr)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(http.StatusNotFound, "no rates for this currency")
	} else if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(cr)
}

// Update
func (h *CurrencyRateHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	existing, err := h.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(http.StatusNotFound, "currency rate not found")
	} else if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	var body struct {
		Value     *float64 `json:"value"`
		Source    *string  `json:"source"`
		FetchedAt *string  `json:"fetched_at"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid JSON")
	}

	if body.Value != nil && *body.Value > 0 {
		existing.Value = *body.Value
	}
	if body.Source != nil {
		src := db.RateSource(*body.Source)
		if src != db.Cifra && src != db.FF && src != db.Manual {
			return fiber.NewError(http.StatusBadRequest, "unsupported source")
		}
		existing.Source = src
	}
	if body.FetchedAt != nil {
		t, err := time.Parse(time.RFC3339, *body.FetchedAt)
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, "invalid fetched_at")
		}
		existing.FetchedAt = t
	}
	existing.UpdatedAt = time.Now().UTC()

	if err := h.repo.Update(existing); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(existing)
}

// Delete
func (h *CurrencyRateHandler) Delete(c *fiber.Ctx) error {
	if err := h.repo.Delete(c.Params("id")); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}
