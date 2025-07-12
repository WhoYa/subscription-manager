package handlers

import (
	"errors"
	"log"
	"strconv"

	repo "github.com/WhoYa/subscription-manager/internal/repository/user"
	dbpkg "github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Healthz(c *fiber.Ctx) error { return c.SendString("ok") }

type UserHandler struct {
	repo repo.UserRepository
}

func NewUserHandler(r repo.UserRepository) *UserHandler {
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
		log.Printf("USER: Failed to parse request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	log.Printf("USER: Creating user request - TGID: %d, Username: %s, Fullname: %s, IsAdmin: %t",
		body.TGID, body.Username, body.Fullname, body.IsAdmin)

	user := dbpkg.User{
		TGID:     body.TGID,
		Username: body.Username,
		Fullname: body.Fullname,
		IsAdmin:  body.IsAdmin,
	}

	log.Printf("USER: Created user struct: %+v", user)

	if err := h.repo.Create(&user); err != nil {
		log.Printf("USER: Failed to create user in database: %v", err)
		// репозиторий уже переводит PG-ошибку дублирования в ErrDuplicateTGID
		if errors.Is(err, repo.ErrDuplicateTGID) {
			return c.Status(409).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("USER: User created successfully: %+v", user)
	return c.Status(201).JSON(user)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	u, err := h.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(u)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	// дефолты
	limit, err := strconv.Atoi(c.Query("limit", "25"))
	if err != nil || limit <= 0 {
		limit = 25
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := h.repo.List(limit, offset)
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

func (h *UserHandler) FindByTGID(c *fiber.Ctx) error {
	tgidStr := c.Params("tgid")
	tgid, err := strconv.ParseInt(tgidStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid tg_id"})
	}

	u, err := h.repo.FindByTGID(tgid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(u)
}
