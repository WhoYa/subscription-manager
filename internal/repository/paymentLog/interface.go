package paymentlog

import (
	"time"

	"github.com/WhoYa/subscription-manager/pkg/db"
)

type PaymentLogRepository interface {
	Create(pl *db.PaymentLog) error
	FindByID(id string) (*db.PaymentLog, error)
	FindByUser(userID string, from, to time.Time) ([]db.PaymentLog, error)
	FindBySubscription(subID string, from, to time.Time) ([]db.PaymentLog, error)
	FindAll(from, to time.Time) ([]db.PaymentLog, error)
}
