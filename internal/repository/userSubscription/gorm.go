package usersubscription

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"gorm.io/gorm"
)

type userSubscriptionGormRepo struct {
	orm *gorm.DB
}

func NewUserSubscriptionRepo(db *gorm.DB) userSubsciptionRepository {
	return &userSubscriptionGormRepo{orm: db}
}

func (r *userSubscriptionGormRepo) Create(us *db.UserSubscription) error {
	return r.orm.Create(us).Error
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

func (r *userSubscriptionGormRepo) FindByUser(userID string) ([]db.UserSubscription, error) {
	var list []db.UserSubscription
	err := r.orm.
		Preload("Subscription").
		Where("user_id = ?", userID).
		Find(&list).Error
	return list, err
}

func (r *userSubscriptionGormRepo) FindBySubscription(subID string) ([]db.UserSubscription, error) {
	var list []db.UserSubscription
	err := r.orm.
		Preload("User").
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
