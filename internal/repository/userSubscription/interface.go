package usersubscription

import "github.com/WhoYa/subscription-manager/pkg/db"

type userSubsciptionRepository interface {
	Create(us *db.UserSubscription) error
	FindByID(id string) (*db.UserSubscription, error)
	FindByUser(userID string) ([]db.UserSubscription, error)
	FindBySubscription(subID string) ([]db.UserSubscription, error)
	UpdateSettings(us *db.UserSubscription) error
	Delete(id string) error
}
