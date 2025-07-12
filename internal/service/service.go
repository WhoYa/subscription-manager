package service

import (
	"time"

	"github.com/WhoYa/subscription-manager/pkg/db"
)

// PaymentAmount представляет рассчитанную сумму к оплате
type PaymentAmount struct {
	UserID         string      `json:"user_id"`
	SubscriptionID string      `json:"subscription_id"`
	Amount         int64       `json:"amount_kopecks"` // копейки для точности
	AmountRubles   float64     `json:"amount_rubles"`  // рубли для удобства
	BaseAmount     float64     `json:"base_amount"`    // "чистая" сумма в рублях
	ProfitAmount   float64     `json:"profit_amount"`  // прибыль в рублях
	Currency       db.Currency `json:"currency"`       // всегда RUB
	ExchangeRate   float64     `json:"exchange_rate"`  // курс конвертации
	DueDate        time.Time   `json:"due_date"`       // дата списания
}

// CurrencyRate представляет курс валюты
type CurrencyRate struct {
	Currency db.Currency `json:"currency"`
	Rate     float64     `json:"rate"`   // курс к рублю
	Source   string      `json:"source"` // источник (Cifra, FF)
}

// ProfitStats представляет статистику прибыли
type ProfitStats struct {
	TotalProfit   float64 `json:"total_profit"`   // общая прибыль в рублях
	TotalPayments int64   `json:"total_payments"` // количество платежей
	AverageProfit float64 `json:"average_profit"` // средняя прибыль за платеж
	Period        string  `json:"period"`         // период (например, "2024-07")
}

// UserProfitStats представляет статистику прибыли по пользователю
type UserProfitStats struct {
	UserID       string  `json:"user_id"`
	Username     string  `json:"username"`
	TotalProfit  float64 `json:"total_profit"`
	PaymentCount int64   `json:"payment_count"`
}

// SubscriptionProfitStats представляет статистику прибыли по подписке
type SubscriptionProfitStats struct {
	SubscriptionID string  `json:"subscription_id"`
	ServiceName    string  `json:"service_name"`
	TotalProfit    float64 `json:"total_profit"`
	PaymentCount   int64   `json:"payment_count"`
}

// Service интерфейс для основной бизнес-логики
type Service interface {
	// CalculateUserPayment рассчитывает сумму к оплате для пользователя по подписке
	// за сутки до даты списания
	CalculateUserPayment(userID, subscriptionID string, dueDate time.Time) (*PaymentAmount, error)
}

// ProfitAnalytics интерфейс для аналитики прибыли (только для администраторов)
type ProfitAnalytics interface {
	// GetMonthlyProfit возвращает общую прибыль за месяц
	GetMonthlyProfit(year int, month int) (*ProfitStats, error)

	// GetUserProfitStats возвращает статистику прибыли по пользователям за период
	GetUserProfitStats(from, to time.Time) ([]UserProfitStats, error)

	// GetSubscriptionProfitStats возвращает статистику прибыли по подпискам за период
	GetSubscriptionProfitStats(from, to time.Time) ([]SubscriptionProfitStats, error)

	// GetTotalProfit возвращает общую прибыль за все время
	GetTotalProfit() (*ProfitStats, error)
}
