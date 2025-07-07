package subscription

import (
	"github.com/WhoYa/subscription-manager/pkg/db"
	"gorm.io/gorm"
)

type subscriptionGormRepo struct {
	orm *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) SubscriptionRepository {
	return &subscriptionGormRepo{orm: db}
}

func (r *subscriptionGormRepo) Create(s *db.Subscription) error {
	return r.orm.Create(s).Error
}

func (r *subscriptionGormRepo) List(limit, offset int) ([]db.Subscription, error) {
	var subscriptions []db.Subscription
	err := r.orm.
		Preload("Users").
		Limit(limit).
		Offset(offset).
		Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionGormRepo) FindByID(id string) (*db.Subscription, error) {
	var s db.Subscription
	err := r.orm.
		Preload("Users").
		First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &s, err
}

func (r *subscriptionGormRepo) Update(s *db.Subscription) error {
	return r.orm.Save(s).Error
}

func (r *subscriptionGormRepo) Delete(id string) error {
	return r.orm.Delete(&db.Subscription{}, "id = ?", id).Error
}
