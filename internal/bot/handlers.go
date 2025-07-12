package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/WhoYa/subscription-manager/internal/bot/api"
	"github.com/WhoYa/subscription-manager/internal/bot/keyboards"
	"github.com/WhoYa/subscription-manager/internal/bot/types"
)

// Subscription handlers

// startCreateSubscription начинает процесс создания подписки
func (b *Bot) startCreateSubscription(userID, chatID int64) {
	userState := b.getUserState(userID)
	userState.State = types.StateAwaitingSubscriptionName
	userState.SubscriptionData = &types.SubscriptionCreateData{}

	b.sendMessage(chatID, "📝 Создание новой подписки\n\nВведите название сервиса (например: Netflix, Spotify):")
}

// handleSubscriptionNameInput обрабатывает ввод названия подписки
func (b *Bot) handleSubscriptionNameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	serviceName := strings.TrimSpace(message.Text)
	if serviceName == "" {
		b.sendMessage(message.Chat.ID, "❌ Название сервиса не может быть пустым. Попробуйте снова:")
		return
	}

	userState.SubscriptionData.ServiceName = serviceName
	userState.State = types.StateAwaitingSubscriptionPrice

	b.sendMessage(message.Chat.ID, fmt.Sprintf("✅ Название: %s\n\nТеперь введите базовую стоимость (например: 9.99):", serviceName))
}

// handleSubscriptionPriceInput обрабатывает ввод цены подписки
func (b *Bot) handleSubscriptionPriceInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	priceStr := strings.TrimSpace(message.Text)
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		b.sendMessage(message.Chat.ID, "❌ Неверный формат цены. Введите число больше 0 (например: 9.99):")
		return
	}

	userState.SubscriptionData.BasePrice = price
	userState.State = types.StateAwaitingSubscriptionCurrency

	text := fmt.Sprintf("✅ Цена: %.2f\n\nВыберите валюту:", price)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.CurrencyKeyboard()
	b.API.Send(msg)
}

// handleCurrencySelection обрабатывает выбор валюты
func (b *Bot) handleCurrencySelection(query *tgbotapi.CallbackQuery) {
	userState := b.getUserState(query.From.ID)

	if userState.State != types.StateAwaitingSubscriptionCurrency {
		b.sendMessage(query.Message.Chat.ID, "❌ Неожиданное действие.")
		return
	}

	currency := strings.TrimPrefix(query.Data, "currency_")
	userState.SubscriptionData.BaseCurrency = currency
	userState.State = types.StateAwaitingSubscriptionPeriod

	text := fmt.Sprintf("✅ Валюта: %s\n\nВведите период списания в днях (например: 30 для ежемесячной подписки):", currency)
	b.sendMessage(query.Message.Chat.ID, text)
}

// handleSubscriptionPeriodInput обрабатывает ввод периода подписки
func (b *Bot) handleSubscriptionPeriodInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	periodStr := strings.TrimSpace(message.Text)
	period, err := strconv.Atoi(periodStr)
	if err != nil || period <= 0 {
		b.sendMessage(message.Chat.ID, "❌ Неверный формат периода. Введите целое число больше 0:")
		return
	}

	userState.SubscriptionData.PeriodDays = period

	// Показываем итоговую информацию и просим подтверждения
	data := userState.SubscriptionData
	text := fmt.Sprintf(`
📝 Подтверждение создания подписки:

🏷️ Сервис: %s
💰 Цена: %.2f %s
📅 Период: %d дней

Создать подписку?`, data.ServiceName, data.BasePrice, data.BaseCurrency, data.PeriodDays)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.ConfirmKeyboard("create_subscription")
	b.API.Send(msg)
}

// User handlers

// startCreateUser начинает процесс создания пользователя
func (b *Bot) startCreateUser(userID, chatID int64) {
	userState := b.getUserState(userID)
	userState.State = types.StateAwaitingUserFullname
	userState.UserCreateData = &types.UserCreateData{}

	b.sendMessage(chatID, "👤 Создание нового пользователя\n\nВведите ФИО пользователя:")
}

// handleUserFullnameInput обрабатывает ввод ФИО пользователя
func (b *Bot) handleUserFullnameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	fullname := strings.TrimSpace(message.Text)
	if fullname == "" {
		b.sendMessage(message.Chat.ID, "❌ ФИО не может быть пустым. Попробуйте снова:")
		return
	}

	userState.UserCreateData.Fullname = fullname
	userState.State = types.StateAwaitingUserTGID

	b.sendMessage(message.Chat.ID, fmt.Sprintf("✅ ФИО: %s\n\nВведите Telegram ID пользователя (числовой ID):", fullname))
}

// handleUserTGIDInput обрабатывает ввод Telegram ID пользователя
func (b *Bot) handleUserTGIDInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	tgidStr := strings.TrimSpace(message.Text)
	tgid, err := strconv.ParseInt(tgidStr, 10, 64)
	if err != nil || tgid <= 0 {
		b.sendMessage(message.Chat.ID, "❌ Неверный формат Telegram ID. Введите положительное число:")
		return
	}

	userState.UserCreateData.TGID = tgid
	userState.State = types.StateAwaitingUserUsername

	b.sendMessage(message.Chat.ID, fmt.Sprintf("✅ Telegram ID: %d\n\nВведите username пользователя (без @, можно оставить пустым):", tgid))
}

// handleUserUsernameInput обрабатывает ввод username пользователя
func (b *Bot) handleUserUsernameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	username := strings.TrimSpace(message.Text)
	// Username может быть пустым
	userState.UserCreateData.Username = username

	// Показываем итоговую информацию и просим подтверждения
	data := userState.UserCreateData
	usernameText := "не указан"
	if username != "" {
		usernameText = "@" + username
	}

	text := fmt.Sprintf(`
👤 Подтверждение создания пользователя:

👤 ФИО: %s
🆔 Telegram ID: %d
📝 Username: %s

Создать пользователя?`, data.Fullname, data.TGID, usernameText)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.ConfirmKeyboard("create_user")
	b.API.Send(msg)
}

// Global settings handlers

// startEditGlobalMarkup начинает процесс редактирования глобальной надбавки
func (b *Bot) startEditGlobalMarkup(userID, chatID int64) {
	userState := b.getUserState(userID)
	userState.State = types.StateAwaitingGlobalMarkup

	// Получаем текущие настройки для отображения
	settings, err := b.Context.APIClient.GetGlobalSettings()
	currentMarkup := "неизвестна"
	if err == nil {
		currentMarkup = fmt.Sprintf("%.2f%%", settings.GlobalMarkupPercent)
	}

	text := fmt.Sprintf(`
📝 Изменение глобальной надбавки

Текущая надбавка: %s

Введите новое значение надбавки в процентах (например: 15.5):`, currentMarkup)

	b.sendMessage(chatID, text)
}

// handleGlobalMarkupInput обрабатывает ввод глобальной надбавки
func (b *Bot) handleGlobalMarkupInput(message *tgbotapi.Message) {
	markupStr := strings.TrimSpace(message.Text)
	markup, err := strconv.ParseFloat(markupStr, 64)
	if err != nil || markup < 0 {
		b.sendMessage(message.Chat.ID, "❌ Неверный формат надбавки. Введите число больше или равное 0:")
		return
	}

	// Обновляем глобальные настройки через API
	req := api.UpdateGlobalSettingsRequest{
		GlobalMarkupPercent: markup,
	}

	settings, err := b.Context.APIClient.UpdateGlobalSettings(req)
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Ошибка при сохранении настроек: %v", err))
		return
	}

	b.setUserState(message.From.ID, types.StateIdle)

	text := fmt.Sprintf("✅ Глобальная надбавка установлена: %.2f%%", settings.GlobalMarkupPercent)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard("global_settings")
	b.API.Send(msg)
}

// General handlers

// handleConfirmation обрабатывает подтверждения действий
func (b *Bot) handleConfirmation(query *tgbotapi.CallbackQuery) {
	action := strings.TrimPrefix(query.Data, "confirm_")
	userState := b.getUserState(query.From.ID)

	switch action {
	case "create_subscription":
		b.confirmCreateSubscription(query.From.ID, query.Message.Chat.ID, userState)
	case "create_user":
		b.confirmCreateUser(query.From.ID, query.Message.Chat.ID, userState)
	default:
		b.sendMessage(query.Message.Chat.ID, "❌ Неизвестное действие.")
	}
}

// confirmCreateSubscription подтверждает создание подписки
func (b *Bot) confirmCreateSubscription(userID, chatID int64, userState *types.UserData) {
	if userState.SubscriptionData == nil {
		b.sendMessage(chatID, "❌ Данные для создания подписки не найдены.")
		return
	}

	data := userState.SubscriptionData

	// Создаем подписку через API
	req := api.CreateSubscriptionRequest{
		ServiceName:  data.ServiceName,
		BasePrice:    data.BasePrice,
		BaseCurrency: data.BaseCurrency,
		PeriodDays:   data.PeriodDays,
	}

	subscription, err := b.Context.APIClient.CreateSubscription(req)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при создании подписки: %v", err))
		return
	}

	text := fmt.Sprintf("✅ Подписка '%s' успешно создана!\n\n💰 Цена: %.2f %s\n📅 Период: %d дней\n🆔 ID: %s",
		subscription.ServiceName, subscription.BasePrice, subscription.BaseCurrency, subscription.PeriodDays, subscription.ID)

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.SubscriptionData = nil

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard("manage_subscriptions")
	b.API.Send(msg)
}

// confirmCreateUser подтверждает создание пользователя
func (b *Bot) confirmCreateUser(userID, chatID int64, userState *types.UserData) {
	if userState.UserCreateData == nil {
		b.sendMessage(chatID, "❌ Данные для создания пользователя не найдены.")
		return
	}

	data := userState.UserCreateData

	// Создаем пользователя через API
	req := api.CreateUserRequest{
		TGID:     data.TGID,
		Username: data.Username,
		Fullname: data.Fullname,
		IsAdmin:  false, // По умолчанию не админ
	}

	user, err := b.Context.APIClient.CreateUser(req)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при создании пользователя: %v", err))
		return
	}

	usernameText := "не указан"
	if user.Username != "" {
		usernameText = "@" + user.Username
	}

	text := fmt.Sprintf("✅ Пользователь успешно создан!\n\n👤 ФИО: %s\n🆔 Telegram ID: %d\n📝 Username: %s\n🆔 ID: %s",
		user.Fullname, user.TGID, usernameText, user.ID)

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.UserCreateData = nil

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard("manage_users")
	b.API.Send(msg)
}

// cancelCurrentOperation отменяет текущую операцию
func (b *Bot) cancelCurrentOperation(userID, chatID int64) {
	userState := b.getUserState(userID)
	userState.State = types.StateIdle
	userState.SubscriptionData = nil
	userState.UserCreateData = nil
	userState.CurrentEntityID = ""

	b.sendMessage(chatID, "❌ Операция отменена.")
	b.showMainMenu(chatID)
}
