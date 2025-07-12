package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	crRepo "github.com/WhoYa/subscription-manager/internal/repository/currencyrate"
	gsRepo "github.com/WhoYa/subscription-manager/internal/repository/globalsettings"
	subRepo "github.com/WhoYa/subscription-manager/internal/repository/subscription"
	usRepo "github.com/WhoYa/subscription-manager/internal/repository/usersubscription"
	"github.com/WhoYa/subscription-manager/pkg/db"
	"gorm.io/gorm"
)

var (
	ErrUserSubscriptionNotFound = errors.New("user subscription not found")
	ErrSubscriptionNotFound     = errors.New("subscription not found")
	ErrExchangeRateNotFound     = errors.New("exchange rate not found")
)

// paymentService простая реализация Service
type paymentService struct {
	userSubRepo  usRepo.UserSubscriptionRepository
	subRepo      subRepo.SubscriptionRepository
	currencyRepo crRepo.CurrencyRateRepository
	settingsRepo gsRepo.GlobalSettingsRepository
}

// NewService создаёт новый экземпляр сервиса
func NewService(
	userSubRepo usRepo.UserSubscriptionRepository,
	subRepo subRepo.SubscriptionRepository,
	currencyRepo crRepo.CurrencyRateRepository,
	settingsRepo gsRepo.GlobalSettingsRepository,
) Service {
	return &paymentService{
		userSubRepo:  userSubRepo,
		subRepo:      subRepo,
		currencyRepo: currencyRepo,
		settingsRepo: settingsRepo,
	}
}

// CalculateUserPayment рассчитывает сумму к оплате для пользователя
func (s *paymentService) CalculateUserPayment(userID, subscriptionID string, dueDate time.Time) (*PaymentAmount, error) {
	// Получаем настройки пользователя для подписки
	userSubs, err := s.userSubRepo.FindByUser(userID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	// Ищем нужную подписку
	var userSub *db.UserSubscription
	for _, us := range userSubs {
		if us.SubscriptionID == subscriptionID {
			userSub = &us
			break
		}
	}

	if userSub == nil {
		return nil, fmt.Errorf("%w: user %s not subscribed to %s", ErrUserSubscriptionNotFound, userID, subscriptionID)
	}

	// Получаем данные подписки
	subscription, err := s.subRepo.FindByID(subscriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: subscription %s", ErrSubscriptionNotFound, subscriptionID)
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Рассчитываем базовую цену
	basePrice := subscription.BasePrice
	baseCurrency := subscription.BaseCurrency

	// Получаем курс валюты (если нужна конвертация)
	exchangeRate := 1.0
	if baseCurrency != db.RUB {
		rate, err := s.currencyRepo.LatestByCurrency(baseCurrency)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("%w: no rate for %s", ErrExchangeRateNotFound, baseCurrency)
			}
			return nil, fmt.Errorf("failed to get exchange rate for %s: %w", baseCurrency, err)
		}
		exchangeRate = rate.Value
	}

	// Конвертируем базовую цену в рубли (это "чистая" сумма)
	baseAmountRub := basePrice * exchangeRate

	// Применяем пользовательские настройки цены
	finalPrice := s.applyPricingMode(baseAmountRub, userSub)

	// Применяем глобальную надбавку, если нет пользовательских настроек
	if userSub.PricingMode == db.None {
		finalPrice = s.applyGlobalMarkup(finalPrice)
	}

	// Вычисляем прибыль
	profitAmount := finalPrice - baseAmountRub

	// Конвертируем в копейки для точности
	amountKopecks := int64(math.Round(finalPrice * 100))

	// Округляем рубли до 2 знаков после запятой
	amountRubles := math.Round(finalPrice*100) / 100
	baseAmountRubles := math.Round(baseAmountRub*100) / 100
	profitAmountRubles := math.Round(profitAmount*100) / 100

	return &PaymentAmount{
		UserID:         userID,
		SubscriptionID: subscriptionID,
		Amount:         amountKopecks,
		AmountRubles:   amountRubles,
		BaseAmount:     baseAmountRubles,
		ProfitAmount:   profitAmountRubles,
		Currency:       db.RUB,
		ExchangeRate:   exchangeRate,
		DueDate:        dueDate,
	}, nil
}

// applyPricingMode применяет пользовательские настройки цены
func (s *paymentService) applyPricingMode(basePrice float64, userSub *db.UserSubscription) float64 {
	switch userSub.PricingMode {
	case db.Percent:
		return basePrice * (1 + userSub.MarkupPercent/100)
	case db.Fixed:
		return userSub.FixedFee
	default: // db.None
		return basePrice
	}
}

// applyGlobalMarkup применяет глобальную надбавку
func (s *paymentService) applyGlobalMarkup(price float64) float64 {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil {
		return price // Нет глобальных настроек
	}

	if settings.GlobalMarkupPercent > 0 {
		return price * (1 + settings.GlobalMarkupPercent/100)
	}

	return price
}
