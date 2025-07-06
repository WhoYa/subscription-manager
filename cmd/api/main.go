package main

import (
	"log"
	"os"

	"github.com/WhoYa/subscription-manager/internal/app"
)

func main() {

	a := app.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := a.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
