package paymentlog

import (
	"time"

	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentLogGormRepo struct {
	orm *gorm.DB
}

func NewPaymentLogRepo(db *gorm.DB) PaymentLogRepository {
	return &paymentLogGormRepo{orm: db}
}

func (r *paymentLogGormRepo) Create(us *db.PaymentLog) error {
	// Генерируем UUID если он не установлен
	if us.ID == "" {
		us.ID = uuid.New().String()
	}

	return r.orm.Create(us).Error
}

func (r *paymentLogGormRepo) FindByID(id string) (*db.PaymentLog, error) {
	var pl db.PaymentLog
	err := r.orm.
		Preload("User").
		Preload("Subscription").
		First(&pl, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &pl, err
}

func (r *paymentLogGormRepo) FindByUser(userID string, from, to time.Time) ([]db.PaymentLog, error) {
	var logs []db.PaymentLog
	err := r.orm.
		Preload("Subscription").
		Where("user_id = ? AND paid_at BETWEEN ? AND ?", userID, from, to).
		Find(&logs).Error
	return logs, err
}

func (r *paymentLogGormRepo) FindBySubscription(subID string, from, to time.Time) ([]db.PaymentLog, error) {
	var logs []db.PaymentLog
	err := r.orm.
		Preload("User").
		Where("subscription_id = ? AND paid_at BETWEEN ? AND ?", subID, from, to).
		Find(&logs).Error
	return logs, err
}

func (r *paymentLogGormRepo) FindAll(from, to time.Time) ([]db.PaymentLog, error) {
	var logs []db.PaymentLog
	err := r.orm.
		Preload("User").
		Preload("Subscription").
		Where("paid_at BETWEEN ? AND ?", from, to).
		Find(&logs).Error
	return logs, err
}
