package handlers

import (
	"github.com/WhoYa/subscription-manager/internal/repository/user"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
)

func Healthz(c *fiber.Ctx) error { return c.SendString("ok") }

type UserHandler struct {
	repo user.UserRepository
}

func NewUserHandler(r user.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var body struct {
		TGID     int64  `json:"tg_id"`
		Username string `json:"username"`
		Fullname string `json:"fullname"`
		IsAdmin  bool   `json:"is_admin"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	user := db.User{
		TGID:     body.TGID,
		Username: body.Username,
		Fullname: body.Fullname,
		IsAdmin:  body.IsAdmin,
	}
	if err := h.repo.Create(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(user)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.repo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error})
	}
	return c.JSON(user)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.repo.List(25, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	var body struct {
		Username *string `json:"username"`
		Fullname *string `json:"fullname"`
		IsAdmin  *bool   `json:"is_admin"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if body.Username != nil {
		user.Username = *body.Username
	}
	if body.Fullname != nil {
		user.Fullname = *body.Fullname
	}
	if body.IsAdmin != nil {
		user.IsAdmin = *body.IsAdmin
	}
	if err := h.repo.Update(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
