package main

import (
	"flag"
	"log"
	"os"

	"github.com/WhoYa/subscription-manager/internal/app"
	"github.com/WhoYa/subscription-manager/internal/util/healthcheck"
)

var healthCheck = flag.Bool("health", false, "run health check and exit")

func main() {
	flag.Parse()

	if *healthCheck {
		if err := healthcheck.Run(); err != nil {
			log.Fatalf("Health check failed: %v", err)
		}
		os.Exit(0)
	}

	a := app.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := a.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
