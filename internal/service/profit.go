package service

import (
	"fmt"
	"time"

	payRepo "github.com/WhoYa/subscription-manager/internal/repository/paymentlog"
	subRepo "github.com/WhoYa/subscription-manager/internal/repository/subscription"
	userRepo "github.com/WhoYa/subscription-manager/internal/repository/user"
)

// profitAnalytics реализация ProfitAnalytics
type profitAnalytics struct {
	paymentRepo payRepo.PaymentLogRepository
	userRepo    userRepo.UserRepository
	subRepo     subRepo.SubscriptionRepository
}

// NewProfitAnalytics создает новый экземпляр сервиса аналитики прибыли
func NewProfitAnalytics(
	paymentRepo payRepo.PaymentLogRepository,
	userRepo userRepo.UserRepository,
	subRepo subRepo.SubscriptionRepository,
) ProfitAnalytics {
	return &profitAnalytics{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
		subRepo:     subRepo,
	}
}

// GetMonthlyProfit возвращает общую прибыль за месяц
func (p *profitAnalytics) GetMonthlyProfit(year int, month int) (*ProfitStats, error) {
	// Определяем границы месяца
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0).Add(-time.Second) // последняя секунда месяца

	return p.calculateProfitForPeriod(from, to, fmt.Sprintf("%04d-%02d", year, month))
}

// GetUserProfitStats возвращает статистику прибыли по пользователям за период
func (p *profitAnalytics) GetUserProfitStats(from, to time.Time) ([]UserProfitStats, error) {
	// Получаем все платежи за период
	payments, err := p.paymentRepo.FindAll(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	// Группируем по пользователям
	userStats := make(map[string]*UserProfitStats)

	for _, payment := range payments {
		userID := payment.UserID

		if _, exists := userStats[userID]; !exists {
			// Получаем данные пользователя
			user, err := p.userRepo.FindByID(userID)
			if err != nil {
				continue // Пропускаем если пользователь не найден
			}

			userStats[userID] = &UserProfitStats{
				UserID:       userID,
				Username:     user.Username,
				TotalProfit:  0,
				PaymentCount: 0,
			}
		}

		// Прибыль в рублях = копейки / 100
		profitRubles := float64(payment.ProfitAmount) / 100
		userStats[userID].TotalProfit += profitRubles
		userStats[userID].PaymentCount++
	}

	// Конвертируем в слайс
	result := make([]UserProfitStats, 0, len(userStats))
	for _, stats := range userStats {
		result = append(result, *stats)
	}

	return result, nil
}

// GetSubscriptionProfitStats возвращает статистику прибыли по подпискам за период
func (p *profitAnalytics) GetSubscriptionProfitStats(from, to time.Time) ([]SubscriptionProfitStats, error) {
	// Получаем все платежи за период
	payments, err := p.paymentRepo.FindAll(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	// Группируем по подпискам
	subStats := make(map[string]*SubscriptionProfitStats)

	for _, payment := range payments {
		subID := payment.SubscriptionID

		if _, exists := subStats[subID]; !exists {
			// Получаем данные подписки
			sub, err := p.subRepo.FindByID(subID)
			if err != nil {
				continue // Пропускаем если подписка не найдена
			}

			subStats[subID] = &SubscriptionProfitStats{
				SubscriptionID: subID,
				ServiceName:    sub.ServiceName,
				TotalProfit:    0,
				PaymentCount:   0,
			}
		}

		// Прибыль в рублях = копейки / 100
		profitRubles := float64(payment.ProfitAmount) / 100
		subStats[subID].TotalProfit += profitRubles
		subStats[subID].PaymentCount++
	}

	// Конвертируем в слайс
	result := make([]SubscriptionProfitStats, 0, len(subStats))
	for _, stats := range subStats {
		result = append(result, *stats)
	}

	return result, nil
}

// GetTotalProfit возвращает общую прибыль за все время
func (p *profitAnalytics) GetTotalProfit() (*ProfitStats, error) {
	// Используем очень широкий диапазон дат для "всего времени"
	from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Now().AddDate(1, 0, 0) // год в будущее

	return p.calculateProfitForPeriod(from, to, "all-time")
}

// calculateProfitForPeriod вспомогательная функция для расчета прибыли за период
func (p *profitAnalytics) calculateProfitForPeriod(from, to time.Time, period string) (*ProfitStats, error) {
	payments, err := p.paymentRepo.FindAll(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments for period: %w", err)
	}

	var totalProfitKopecks int64
	var paymentCount int64

	for _, payment := range payments {
		totalProfitKopecks += payment.ProfitAmount
		paymentCount++
	}

	totalProfitRubles := float64(totalProfitKopecks) / 100

	var averageProfit float64
	if paymentCount > 0 {
		averageProfit = totalProfitRubles / float64(paymentCount)
	}

	return &ProfitStats{
		TotalProfit:   totalProfitRubles,
		TotalPayments: paymentCount,
		AverageProfit: averageProfit,
		Period:        period,
	}, nil
}
