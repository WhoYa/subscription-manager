package subscription

import (
	"errors"
	"strings"

	"github.com/WhoYa/subscription-manager/pkg/db"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	// ErrDuplicateServiceName возвращается, когда уже есть подписка с таким service_name
	ErrDuplicateServiceName = errors.New("duplicate service_name")
)

type subscriptionGormRepo struct {
	orm *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) SubscriptionRepository {
	return &subscriptionGormRepo{orm: db}
}

func (r *subscriptionGormRepo) Create(s *db.Subscription) error {
	err := r.orm.Create(s).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// проверяем имя constraint для точности
			if strings.Contains(pgErr.ConstraintName, "subscriptions_service_name_key") {
				return ErrDuplicateServiceName
			}
			return ErrDuplicateServiceName
		}
	}
	return err
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

func (r *subscriptionGormRepo) FindByServiceName(name string) (*db.Subscription, error) {
	var s db.Subscription
	err := r.orm.
		First(&s, "service_name = ?", name).
		Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionGormRepo) Update(s *db.Subscription) error {
	return r.orm.Save(s).Error
}

func (r *subscriptionGormRepo) Delete(id string) error {
	return r.orm.Delete(&db.Subscription{}, "id = ?", id).Error
}
