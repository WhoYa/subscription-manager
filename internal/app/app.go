package app

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/WhoYa/subscription-manager/internal/handlers"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/WhoYa/subscription-manager/pkg/db/migrations"
)

func New() *fiber.App {

	// DB + Migrations
	gormDB, err := db.Open()
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	m := gormigrate.New(gormDB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.InitialMigration(),
		migrations.AddAllTables(),
	})
	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Println("Migrations applied")

	// Fiber + Routes
	app := fiber.New()
	app.Get("/healthz", handlers.Healthz)

	api := app.Group("/api")
	u := api.Group("/users")
	u.Post("/", handlers.CreateUser(gormDB))
	u.Get("/", handlers.ListUser(gormDB))
	u.Get("/:id", handlers.GetUser(gormDB))
	u.Patch("/:id", handlers.UpdateUser(gormDB))
	u.Delete("/:id", handlers.DeleteUser(gormDB))

	s := api.Group("/subscriptions")
	s.Post("/", handlers.CreateSubscription(gormDB))
	s.Get("/", handlers.ListSubscription(gormDB))
	s.Get("/:id", handlers.GetSubscription(gormDB))
	s.Patch("/:id", handlers.UpdateSubscription(gormDB))
	s.Delete("/:id", handlers.DeleteSubscription(gormDB))

	return app
}
