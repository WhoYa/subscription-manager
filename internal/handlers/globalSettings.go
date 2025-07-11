package handlers

import (
	"errors"
	"net/http"
	"time"

	repo "github.com/WhoYa/subscription-manager/internal/repository/globalsettings"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GlobalSettingsHandler struct {
	repo repo.GlobalSettingsRepository
}

func NewGlobalSettingsHandler(r repo.GlobalSettingsRepository) *GlobalSettingsHandler {
	return &GlobalSettingsHandler{repo: r}
}

func (h *GlobalSettingsHandler) Get(c *fiber.Ctx) error {
	gs, err := h.repo.Get()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "global settings not found"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(gs)
}

func (h *GlobalSettingsHandler) Create(c *fiber.Ctx) error {
	var body struct {
		GlobalMarkupPercent float64 `json:"global_markup_percent"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if existing, _ := h.repo.Get(); existing != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "global settings already exist"})
	}

	gs := db.GlobalSettings{
		GlobalMarkupPercent: body.GlobalMarkupPercent,
	}
	if err := h.repo.Create(&gs); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(gs)
}

func (h *GlobalSettingsHandler) Update(c *fiber.Ctx) error {
	var body struct {
		GlobalMarkupPercent float64 `json:"global_markup_percent"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	gs, err := h.repo.Get()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "global settings not found"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	gs.GlobalMarkupPercent = body.GlobalMarkupPercent
	gs.UpdatedAt = time.Now().UTC()

	if err := h.repo.Update(gs); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(gs)
}
