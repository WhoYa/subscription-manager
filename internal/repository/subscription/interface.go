package subscription

import "github.com/WhoYa/subscription-manager/pkg/db"

type SubscriptionRepository interface {
	Create(s *db.Subscription) error
	List(limit, offset int) ([]db.Subscription, error)
	FindByID(id string) (*db.Subscription, error)
	FindByServiceName(name string) (*db.Subscription, error)
	Update(s *db.Subscription) error
	Delete(id string) error
}
