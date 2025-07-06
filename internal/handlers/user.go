package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/WhoYa/subscription-manager/pkg/db"
)

func Healthz(c *fiber.Ctx) error { return c.SendString("ok") }

func CreateUser(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		if err := dbConn.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error:": err.Error()})
		}
		return c.Status(201).JSON(user)
	}
}

func ListUser(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []db.User
		if err := dbConn.Find(&users).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(users)
	}
}

func GetUser(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user db.User
		err := dbConn.
			Preload("Subscriptions").
			Preload("Payments").
			First(&user, "id = ?", id).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{"error": "user not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	}
}

func UpdateUser(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var body struct {
			Username *string `json:"username"`
			Fullname *string `json:"fullname"`
			IsAdmin  *bool   `json:"is_admin"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}
		var user db.User
		if err := dbConn.First(&user, "id = ?", id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
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
		if err := dbConn.Save(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	}
}

func DeleteUser(dbConn *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := dbConn.Delete(&db.User{}, "id = ?", id).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	}
}
