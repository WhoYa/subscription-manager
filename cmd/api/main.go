package main

import (
	"log"
	"os"

	"github.com/WhoYa/subscription-manager/db"
	"github.com/WhoYa/subscription-manager/db/migrations"
	"github.com/gofiber/fiber/v2"

	"github.com/go-gormigrate/gormigrate/v2"
)

func main() {

	gormDB, err := db.Open()
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	m := gormigrate.New(
		gormDB,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			migrations.InitialMigration(),
			migrations.AddAllTables(),
			//migrations.SeedDemoData(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Println("Migrations applied")

	app := fiber.New()
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })
	if err := app.Listen(":" + os.Getenv(("PORT"))); err != nil {
		log.Fatal("Server failed: &s", err)

	}

}
