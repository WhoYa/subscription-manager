package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:uuid;primaryKey"`
	TGID          int64          `gorm:"uniqueIndex;not null"`
	Username      string         `gorm:"size:200"`
	Fullname      string         `gorm:"size:200"`
	IsAdmin       bool           `gorm:"default:false"`
	Subscriptions []Subscription `gorm:"many2many:user_subscriptions"`
	Payments      []PaymentLog
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Subscription struct {
	ID           string   `gorm:"type:uuid;primaryKey"`
	ServiceName  string   `gorm:"size:200"`
	IconURL      string   `gorm:"size:800"`
	BasePrice    float64  `gorm:"type:numeric(12,2)"`
	BaseCurrency Currency `gorm:"type:currency_enum"`
	IsActive     bool     `gorm:"default:true"`
	Users        []User   `gorm:"many2many:user_subscriptions"`
	PeriodDays   int      `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserSubscription struct {
	ID             string      `gorm:"type:uuid;primaryKey"`
	UserID         string      `gorm:"type:uuid;not null;uniqueIndex:user_sub_uq"`
	SubscriptionID string      `gorm:"type:uuid;not null;uniqueIndex:user_sub_uq"`
	PricingMode    PricingMode `gorm:"type:pricing_mode_enum;default:'none'"`
	MarkupPercent  float64     `gorm:"default:0"`
	FixedFee       float64     `gorm:"default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	User         User         `gorm:"foreignkey:UserID;references:ID"`
	Subscription Subscription `gorm:"foreignkey:SubscriptionID;references:ID"`
}

type PaymentLog struct {
	ID             string   `gorm:"type:uuid;primaryKey"`
	UserID         string   `gorm:"type:uuid;not null"`
	SubscriptionID string   `gorm:"type:uuid;not null"`
	Amount         int64    `gorm:"type:bigint"` // итоговая сумма в копейках
	BaseAmount     int64    `gorm:"type:bigint"` // базовая "чистая" сумма в копейках
	ProfitAmount   int64    `gorm:"type:bigint"` // прибыль в копейках
	Currency       Currency `gorm:"type:currency_enum"`
	RateUsed       float64  `gorm:"not null"`
	PaidAt         time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time

	User         User         `gorm:"foreignkey:UserID;references:ID"`
	Subscription Subscription `gorm:"foreignkey:SubscriptionID;references:ID"`
}

type GlobalSettings struct {
	ID                  string         `gorm:"type:uuid;primaryKey" json:"id"`
	GlobalMarkupPercent float64        `gorm:"default:0" json:"global_markup_percent"`
	UpdatedAt           time.Time      `json:"updated_at"`
	CreatedAt           time.Time      `json:"created_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type CurrencyRate struct {
	ID        string     `gorm:"type:uuid;primaryKey"`
	Currency  Currency   `gorm:"type:currency_enum"`
	Value     float64    `gorm:"not null"`
	Source    RateSource `gorm:"type:ratesource_enum"`
	FetchedAt time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
