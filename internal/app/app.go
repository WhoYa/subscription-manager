package app

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/WhoYa/subscription-manager/internal/handlers"
	"github.com/WhoYa/subscription-manager/internal/repository/user"
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

	userRepo := user.NewUserRepo(gormDB)
	userHandler := handlers.NewUserHandler(userRepo)

	app := fiber.New()
	app.Get("/healthz", handlers.Healthz)
	api := app.Group("/api")

	u := api.Group("/users")
	u.Post("/", userHandler.Create)
	u.Get("/", userHandler.List)
	u.Get("/:id", userHandler.Get)
	u.Patch("/:id", userHandler.Update)
	u.Delete("/:id", userHandler.Delete)

	return app
}
