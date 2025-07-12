package types

import (
	"github.com/WhoYa/subscription-manager/internal/bot/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UserState представляет состояние пользователя в боте
type UserState string

const (
	StateIdle                         UserState = "idle"
	StateAwaitingSubscriptionName     UserState = "awaiting_subscription_name"
	StateAwaitingSubscriptionPrice    UserState = "awaiting_subscription_price"
	StateAwaitingSubscriptionCurrency UserState = "awaiting_subscription_currency"
	StateAwaitingSubscriptionPeriod   UserState = "awaiting_subscription_period"
	StateAwaitingUserFullname         UserState = "awaiting_user_fullname"
	StateAwaitingUserTGID             UserState = "awaiting_user_tgid"
	StateAwaitingUserUsername         UserState = "awaiting_user_username"
	StateAwaitingGlobalMarkup         UserState = "awaiting_global_markup"

	// Состояния для редактирования
	StateEditingSubscriptionName     UserState = "editing_subscription_name"
	StateEditingSubscriptionPrice    UserState = "editing_subscription_price"
	StateEditingSubscriptionCurrency UserState = "editing_subscription_currency"
	StateEditingSubscriptionPeriod   UserState = "editing_subscription_period"
	StateEditingUserFullname         UserState = "editing_user_fullname"
	StateEditingUserUsername         UserState = "editing_user_username"
)

// UserData содержит временные данные для создания/редактирования
type UserData struct {
	State              UserState
	SubscriptionData   *SubscriptionCreateData
	UserCreateData     *UserCreateData
	EditData           *EditData
	CurrentEntityID    string // ID редактируемой сущности
	CurrentMessageID   int    // ID текущего сообщения для редактирования
	CurrentChatID      int64  // ID чата для редактирования сообщения
	CurrentMenuContext string // Контекст текущего меню (subscriptions, users, main)
}

// SubscriptionCreateData временные данные для создания подписки
type SubscriptionCreateData struct {
	ServiceName  string
	BasePrice    float64
	BaseCurrency string
	PeriodDays   int
}

// UserCreateData временные данные для создания пользователя
type UserCreateData struct {
	Fullname string
	TGID     int64
	Username string
}

// EditData содержит данные для редактирования
type EditData struct {
	EntityType     string // "user" или "subscription"
	EntityID       string
	OriginalEntity interface{}            // оригинальная сущность для отображения
	UpdatedFields  map[string]interface{} // обновленные поля
}

// BotContext содержит контекст бота и API клиенты
type BotContext struct {
	Bot          *tgbotapi.BotAPI
	APIClient    *api.Client
	APIBaseURL   string
	UserStates   map[int64]*UserData
	AdminUserIDs []int64
}

// CallbackData представляет структурированные callback данные
type CallbackData struct {
	Action string
	ID     string
	Page   string
}
