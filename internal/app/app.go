package app

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/WhoYa/subscription-manager/internal/handlers"
	crRepo "github.com/WhoYa/subscription-manager/internal/repository/currencyrate"
	gsRepo "github.com/WhoYa/subscription-manager/internal/repository/globalsettings"
	payRepo "github.com/WhoYa/subscription-manager/internal/repository/paymentlog"
	subRepo "github.com/WhoYa/subscription-manager/internal/repository/subscription"
	userRepo "github.com/WhoYa/subscription-manager/internal/repository/user"
	usRepo "github.com/WhoYa/subscription-manager/internal/repository/usersubscription"

	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/WhoYa/subscription-manager/pkg/db/migrations"
)

func New() *fiber.App {

	// DB + Migrations ---------------------------------------------------------
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

	// Repositories ------------------------------------------------------------
	uRepo := userRepo.NewUserRepo(gormDB)
	sRepo := subRepo.NewSubscriptionRepo(gormDB)
	usRepo := usRepo.NewUserSubscriptionRepo(gormDB)
	pRepo := payRepo.NewPaymentLogRepo(gormDB)
	gsRepo := gsRepo.NewGlobalSettingsRepository(gormDB)
	crRepo := crRepo.NewCurrencyRateRepo(gormDB)

	// Handlers ----------------------------------------------------------------
	uH := handlers.NewUserHandler(uRepo)
	sH := handlers.NewSubscriptionHandler(sRepo)
	usH := handlers.NewUserSubscriptionHandler(usRepo)
	pH := handlers.NewPaymentLogHandler(pRepo)
	gsH := handlers.NewGlobalSettingsHandler(gsRepo)
	crH := handlers.NewCurrencyRateHandler(crRepo)

	// Fiber + Routes ----------------------------------------------------------
	app := fiber.New()
	api := app.Group("/api")

	// health
	api.Get("/healthz", handlers.Healthz)

	// users
	u := api.Group("/users")
	u.Post("/", uH.Create)
	u.Get("/", uH.List)
	u.Get("/:id", uH.Get)
	u.Patch("/:id", uH.Update)
	u.Delete("/:id", uH.Delete)

	// users -> subscriptions (user-sub join)
	us := u.Group("/:userID/subscriptions")
	us.Post("/", usH.Create)
	us.Get("/", usH.ListByUser)
	us.Patch("/:id", usH.UpdateSettings)
	us.Delete("/:id", usH.Delete)

	// users -> payments
	up := u.Group("/:userID/payments")
	up.Get("/", pH.ListByUser)
	up.Post("/", pH.Create)

	// subscriptions
	s := api.Group("/subscriptions")
	s.Post("/", sH.Create)
	s.Get("/", sH.List)
	s.Get("/:id", sH.Get)
	s.Patch("/:id", sH.Update)
	s.Delete("/:id", sH.Delete)

	// subscriptions -> payments
	sp := s.Group("/:subID/payments")
	sp.Get("/", pH.ListBySubscription)

	// standalone payments list
	api.Get("/payments", pH.ListAll)

	// global settings (singleton)
	settings := api.Group("/settings")
	settings.Get("/", gsH.Get)
	settings.Post("/", gsH.Create)
	settings.Put("/", gsH.Update)

	// currency rates ------------------------------------------------------
	cr := api.Group("/currency_rates")
	cr.Post("/", crH.Create)                // POST   /api/currency_rates
	cr.Get("/", crH.List)                   // GET    /api/currency_rates?limit=&offset=
	cr.Get("/:id", crH.Get)                 // GET    /api/currency_rates/:id
	cr.Get("/latest/:currency", crH.Latest) // GET    /api/currency_rates/latest/USD
	cr.Put("/:id", crH.Update)              // PUT    /api/currency_rates/:id
	cr.Delete("/:id", crH.Delete)           // DELETE /api/currency_

	return app
}
