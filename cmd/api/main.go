package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WhoYa/subscription-manager/db"
	"github.com/WhoYa/subscription-manager/db/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

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
			migrations.SeedDemoData(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Println("Migrations applied")

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed: &s", err)
	}
}
