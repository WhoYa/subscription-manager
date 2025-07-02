package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/WhoYa/subscription-manager/db"
)

func SeedDemoData() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250704_seed_demo",
		Migrate: func(tx *gorm.DB) error {
			// ——— 1) Пользователи ———
			u1 := &db.User{TGID: 11111, Username: "alice", Fullname: "Alice Example"}
			u2 := &db.User{TGID: 22222, Username: "bob", Fullname: "Bob Example"}
			if err := tx.Create(u1).Error; err != nil {
				return err
			}
			if err := tx.Create(u2).Error; err != nil {
				return err
			}

			// ——— 2) Подписки ———
			now := time.Now()
			s1 := &db.Subscription{
				ServiceName:  "Netflix",
				IconURL:      "https://example.com/netflix.png",
				BasePrice:    10.99,
				BaseCurrency: db.USD,
				Period:       now.AddDate(0, 1, 0), // +1 месяц
				IsActive:     true,
			}
			s2 := &db.Subscription{
				ServiceName:  "Spotify",
				IconURL:      "https://example.com/spotify.png",
				BasePrice:    5.49,
				BaseCurrency: db.USD,
				Period:       now.AddDate(0, 1, 0),
				IsActive:     true,
			}
			if err := tx.Create(s1).Error; err != nil {
				return err
			}
			if err := tx.Create(s2).Error; err != nil {
				return err
			}

			// ——— 3) Глобальные настройки ———
			gs := &db.GlobalSettings{GlobalMarkupPercent: 5.0}
			if err := tx.Create(gs).Error; err != nil {
				return err
			}

			// ——— 4) Курсы валют ———
			cr1 := &db.CurrencyRate{
				Currency:  db.USD,
				Value:     92.50,
				Source:    db.Cifra,
				FetchedAt: now,
			}
			cr2 := &db.CurrencyRate{
				Currency:  db.EUR,
				Value:     100.25,
				Source:    db.FF,
				FetchedAt: now,
			}
			if err := tx.Create(cr1).Error; err != nil {
				return err
			}
			if err := tx.Create(cr2).Error; err != nil {
				return err
			}

			// ——— 5) Связи User–Subscription (персональный тариф) ———
			us1 := &db.UserSubscription{
				UserID:         u1.ID,
				SubscriptionID: s1.ID,
				PricingMode:    db.Percent,
				MarkupPercent:  2.5,
			}
			us2 := &db.UserSubscription{
				UserID:         u2.ID,
				SubscriptionID: s2.ID,
				PricingMode:    db.Fixed,
				FixedFee:       2000.0,
			}
			if err := tx.Create(us1).Error; err != nil {
				return err
			}
			if err := tx.Create(us2).Error; err != nil {
				return err
			}

			// ——— 6) Логи платежей ———
			pl1 := &db.PaymentLog{
				UserID:         u1.ID,
				SubscriptionID: s1.ID,
				Amount:         1099,   // в копейках: 10.99×100
				Currency:       db.RUB, // если есть ENUM RUB
				RateUsed:       92.50,
				PaidAt:         now,
			}
			pl2 := &db.PaymentLog{
				UserID:         u2.ID,
				SubscriptionID: s2.ID,
				Amount:         549,
				Currency:       db.RUB,
				RateUsed:       92.50,
				PaidAt:         now,
			}
			if err := tx.Create(pl1).Error; err != nil {
				return err
			}
			if err := tx.Create(pl2).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Удаляем seed-данные по “сырым” SQL
			if err := tx.Exec("DELETE FROM payment_logs").Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_subscriptions").Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM currency_rates").Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM global_settings").Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM subscriptions").Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM users").Error; err != nil {
				return err
			}
			return nil
		},
	}
}
