package usersubscription

import (
	"errors"
	"strings"

	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrDuplicateUserSubscription = errors.New("duplicate user_subscription")

type userSubscriptionGormRepo struct {
	orm *gorm.DB
}

func NewUserSubscriptionRepo(db *gorm.DB) UserSubscriptionRepository {
	return &userSubscriptionGormRepo{orm: db}
}

func (r *userSubscriptionGormRepo) Create(us *db.UserSubscription) error {
	// Генерируем UUID если он не установлен
	if us.ID == "" {
		us.ID = uuid.New().String()
	}

	err := r.orm.Create(us).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// проверяем название уникального индексa,
			// скорее всего usersubscriptions_user_id_subscription_id_key
			if strings.Contains(pgErr.ConstraintName, "user_subscriptions_user_id_subscription_id") {
				return ErrDuplicateUserSubscription
			}
			return ErrDuplicateUserSubscription
		}
	}
	return err
}

func (r *userSubscriptionGormRepo) FindByID(id string) (*db.UserSubscription, error) {
	var us db.UserSubscription
	err := r.orm.
		Preload("User").
		Preload("Subscription").
		First(&us, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &us, err
}

func (r *userSubscriptionGormRepo) FindByUser(userID string, limit, offset int) ([]db.UserSubscription, error) {
	var list []db.UserSubscription
	err := r.orm.
		Preload("User").
		Preload("Subscription").
		Where("user_id = ?", userID).
		Find(&list).Error
	return list, err
}

func (r *userSubscriptionGormRepo) FindBySubscription(subID string) ([]db.UserSubscription, error) {
	var list []db.UserSubscription
	err := r.orm.
		Preload("User").
		Preload("Subscription").
		Where("subscription_id = ?", subID).
		Find(&list).Error
	return list, err
}
func (r *userSubscriptionGormRepo) UpdateSettings(us *db.UserSubscription) error {
	return r.orm.Model(us).Select("PricingMode", "MarkupPercent", "FixedFee").Updates(us).Error
}
func (r *userSubscriptionGormRepo) Delete(id string) error {
	return r.orm.Delete(&db.UserSubscription{}, "id = ?", id).Error
}
