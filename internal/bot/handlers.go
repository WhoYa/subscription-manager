package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/WhoYa/subscription-manager/internal/bot/api"
	"github.com/WhoYa/subscription-manager/internal/bot/keyboards"
	"github.com/WhoYa/subscription-manager/internal/bot/types"
)

// Subscription handlers

// startCreateSubscription начинает процесс создания подписки
func (b *Bot) startCreateSubscription(userID, chatID int64, messageID int) {
	userState := b.getUserState(userID)
	userState.State = types.StateAwaitingSubscriptionName
	userState.SubscriptionData = &types.SubscriptionCreateData{}
	userState.CurrentMenuContext = "subscriptions"
	userState.CurrentMessageID = messageID
	userState.CurrentChatID = chatID

	keyboard := keyboards.CreateProcessKeyboard("start")
	b.editMessage(chatID, messageID, MessageSubscriptionCreateStart, &keyboard)
}

// handleSubscriptionNameInput обрабатывает ввод названия подписки
func (b *Bot) handleSubscriptionNameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	log.Printf("SUBSCRIPTION_NAME: Processing input from user %d: %s", message.From.ID, message.Text)
	log.Printf("SUBSCRIPTION_NAME: Current userState: %+v", userState)

	serviceName, err := validateString(message.Text, false)
	if err != nil {
		log.Printf("SUBSCRIPTION_NAME: Validation failed for user %d: %v", message.From.ID, err)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, MessageSubscriptionNameEmpty, &keyboard)
		return
	}

	if userState.SubscriptionData == nil {
		log.Printf("SUBSCRIPTION_NAME: SubscriptionData is nil for user %d, creating new", message.From.ID)
		userState.SubscriptionData = &types.SubscriptionCreateData{}
	}

	userState.SubscriptionData.ServiceName = serviceName
	userState.State = types.StateAwaitingSubscriptionPrice

	log.Printf("SUBSCRIPTION_NAME: Set service name for user %d: %s", message.From.ID, serviceName)
	log.Printf("SUBSCRIPTION_NAME: Updated SubscriptionData: %+v", userState.SubscriptionData)

	text := fmt.Sprintf(MessageSubscriptionPriceStep, serviceName)
	keyboard := keyboards.CreateProcessKeyboard("process")
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// handleSubscriptionPriceInput обрабатывает ввод цены подписки
func (b *Bot) handleSubscriptionPriceInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	log.Printf("SUBSCRIPTION_PRICE: Processing input from user %d: %s", message.From.ID, message.Text)
	log.Printf("SUBSCRIPTION_PRICE: Current SubscriptionData: %+v", userState.SubscriptionData)

	price, err := validateFloat64(message.Text, 0.01)
	if err != nil {
		log.Printf("SUBSCRIPTION_PRICE: Validation failed for user %d: %v", message.From.ID, err)
		text := fmt.Sprintf(MessageSubscriptionPriceError, userState.SubscriptionData.ServiceName)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
		return
	}

	if userState.SubscriptionData == nil {
		log.Printf("SUBSCRIPTION_PRICE: ERROR - SubscriptionData is nil for user %d", message.From.ID)
		userState.SubscriptionData = &types.SubscriptionCreateData{}
	}

	userState.SubscriptionData.BasePrice = price
	userState.State = types.StateAwaitingSubscriptionCurrency

	log.Printf("SUBSCRIPTION_PRICE: Set price for user %d: %.2f", message.From.ID, price)
	log.Printf("SUBSCRIPTION_PRICE: Updated SubscriptionData: %+v", userState.SubscriptionData)

	text := fmt.Sprintf(MessageSubscriptionCurrencyStep, userState.SubscriptionData.ServiceName, price)
	keyboard := keyboards.CurrencyKeyboardWithNav()
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// handleSubscriptionPeriodInput обрабатывает ввод периода подписки
func (b *Bot) handleSubscriptionPeriodInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	log.Printf("SUBSCRIPTION_PERIOD: Processing input from user %d: %s", message.From.ID, message.Text)
	log.Printf("SUBSCRIPTION_PERIOD: Current SubscriptionData: %+v", userState.SubscriptionData)

	period, err := validateInt(message.Text, 1)
	if err != nil {
		log.Printf("SUBSCRIPTION_PERIOD: Validation failed for user %d: %v", message.From.ID, err)
		text := fmt.Sprintf(MessageSubscriptionPeriodError, userState.SubscriptionData.ServiceName, userState.SubscriptionData.BasePrice, userState.SubscriptionData.BaseCurrency)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
		return
	}

	if userState.SubscriptionData == nil {
		log.Printf("SUBSCRIPTION_PERIOD: ERROR - SubscriptionData is nil for user %d", message.From.ID)
		userState.SubscriptionData = &types.SubscriptionCreateData{}
	}

	userState.SubscriptionData.PeriodDays = period

	log.Printf("SUBSCRIPTION_PERIOD: Set period for user %d: %d days", message.From.ID, period)
	log.Printf("SUBSCRIPTION_PERIOD: Final SubscriptionData: %+v", userState.SubscriptionData)

	// Показываем итоговую информацию и просим подтверждения
	data := userState.SubscriptionData
	text := fmt.Sprintf(MessageSubscriptionConfirm, data.ServiceName, data.BasePrice, data.BaseCurrency, data.PeriodDays)

	keyboard := keyboards.CreateConfirmKeyboard("create_subscription")
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// User handlers

// startCreateUser начинает процесс создания пользователя
func (b *Bot) startCreateUser(userID, chatID int64, messageID int) {
	userState := b.getUserState(userID)
	userState.State = types.StateAwaitingUserFullname
	userState.UserCreateData = &types.UserCreateData{}
	userState.CurrentMenuContext = "users"
	userState.CurrentMessageID = messageID
	userState.CurrentChatID = chatID

	keyboard := keyboards.CreateProcessKeyboard("start")
	b.editMessage(chatID, messageID, MessageUserCreateStart, &keyboard)
}

// handleUserFullnameInput обрабатывает ввод ФИО пользователя
func (b *Bot) handleUserFullnameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	log.Printf("USER_FULLNAME: Processing input from user %d: %s", message.From.ID, message.Text)
	log.Printf("USER_FULLNAME: Current userState: %+v", userState)

	fullname, err := validateString(message.Text, false)
	if err != nil {
		log.Printf("USER_FULLNAME: Validation failed for user %d: %v", message.From.ID, err)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, MessageUserFullnameEmpty, &keyboard)
		return
	}

	if userState.UserCreateData == nil {
		log.Printf("USER_FULLNAME: UserCreateData is nil for user %d, creating new", message.From.ID)
		userState.UserCreateData = &types.UserCreateData{}
	}

	userState.UserCreateData.Fullname = fullname
	userState.State = types.StateAwaitingUserTGID

	log.Printf("USER_FULLNAME: Set fullname for user %d: %s", message.From.ID, fullname)
	log.Printf("USER_FULLNAME: Updated UserCreateData: %+v", userState.UserCreateData)

	text := fmt.Sprintf(MessageUserTGIDStep, fullname)
	keyboard := keyboards.CreateProcessKeyboard("process")
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// handleUserTGIDInput обрабатывает ввод Telegram ID пользователя
func (b *Bot) handleUserTGIDInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	log.Printf("USER_TGID: Processing input from user %d: %s", message.From.ID, message.Text)
	log.Printf("USER_TGID: Current UserCreateData: %+v", userState.UserCreateData)

	tgid, err := validateInt64(message.Text, 1)
	if err != nil {
		log.Printf("USER_TGID: Validation failed for user %d: %v", message.From.ID, err)
		text := fmt.Sprintf(MessageUserTGIDError, userState.UserCreateData.Fullname)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
		return
	}

	if userState.UserCreateData == nil {
		log.Printf("USER_TGID: ERROR - UserCreateData is nil for user %d", message.From.ID)
		userState.UserCreateData = &types.UserCreateData{}
	}

	userState.UserCreateData.TGID = tgid
	userState.State = types.StateAwaitingUserUsername

	log.Printf("USER_TGID: Set TGID for user %d: %d", message.From.ID, tgid)
	log.Printf("USER_TGID: Updated UserCreateData: %+v", userState.UserCreateData)

	text := fmt.Sprintf(MessageUserUsernameStep, userState.UserCreateData.Fullname, tgid)
	keyboard := keyboards.CreateProcessKeyboard("process")
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// handleUserUsernameInput обрабатывает ввод username пользователя
func (b *Bot) handleUserUsernameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	log.Printf("USER_USERNAME: Processing input from user %d: %s", message.From.ID, message.Text)
	log.Printf("USER_USERNAME: Current UserCreateData: %+v", userState.UserCreateData)

	username, _ := validateString(message.Text, true) // Username может быть пустым

	if userState.UserCreateData == nil {
		log.Printf("USER_USERNAME: ERROR - UserCreateData is nil for user %d", message.From.ID)
		userState.UserCreateData = &types.UserCreateData{}
	}

	userState.UserCreateData.Username = username

	log.Printf("USER_USERNAME: Set username for user %d: %s", message.From.ID, username)
	log.Printf("USER_USERNAME: Final UserCreateData: %+v", userState.UserCreateData)

	// Показываем итоговую информацию и просим подтверждения
	data := userState.UserCreateData
	usernameText := formatUsername(username)

	text := fmt.Sprintf(MessageUserConfirm, data.Fullname, data.TGID, usernameText)

	keyboard := keyboards.CreateConfirmKeyboard("create_user")
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// Global settings handlers

// startEditGlobalMarkup начинает процесс редактирования глобальной надбавки
func (b *Bot) startEditGlobalMarkup(userID, chatID int64, messageID int) {
	userState := b.getUserState(userID)
	userState.State = types.StateAwaitingGlobalMarkup
	userState.CurrentMenuContext = "global_settings"
	userState.CurrentMessageID = messageID
	userState.CurrentChatID = chatID

	// Получаем текущие настройки для отображения
	settings, err := b.Context.APIClient.GetGlobalSettings()
	currentMarkup := StatusUnknown
	if err == nil {
		currentMarkup = fmt.Sprintf("%.2f%%", settings.GlobalMarkupPercent)
	}

	text := fmt.Sprintf(MessageGlobalMarkupStart, currentMarkup)

	keyboard := keyboards.CreateProcessKeyboard("start")
	b.editMessage(chatID, messageID, text, &keyboard)
}

// handleGlobalMarkupInput обрабатывает ввод глобальной надбавки
func (b *Bot) handleGlobalMarkupInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)

	markup, err := validateFloat64(message.Text, 0)
	if err != nil {
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, MessageGlobalMarkupError, &keyboard)
		return
	}

	logInfo("GlobalMarkup", fmt.Sprintf("User input: %s, parsed value: %.2f%%", message.Text, markup))

	// Сначала пытаемся обновить настройки
	updateReq := api.UpdateGlobalSettingsRequest{
		GlobalMarkupPercent: markup,
	}

	logInfo("GlobalMarkup", fmt.Sprintf("Sending update request: %+v", updateReq))

	settings, err := b.Context.APIClient.UpdateGlobalSettings(updateReq)
	if err != nil {
		logError("UpdateGlobalSettings", err)

		// Если ошибка 404, значит настройки не существуют, создаем их
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			logInfo("GlobalMarkup", "Settings not found, creating new ones")
			createReq := api.CreateGlobalSettingsRequest{
				GlobalMarkupPercent: markup,
			}

			logInfo("GlobalMarkup", fmt.Sprintf("Sending create request: %+v", createReq))

			settings, err = b.Context.APIClient.CreateGlobalSettings(createReq)
			if err != nil {
				logError("CreateGlobalSettings", err)
				errorText := fmt.Sprintf(MessageError, handleAPIError(err, "CreateGlobalSettings"))
				keyboard := keyboards.CreateSuccessKeyboard("global_settings")
				b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, errorText, &keyboard)
				return
			}
		} else {
			errorText := fmt.Sprintf(MessageError, handleAPIError(err, "UpdateGlobalSettings"))
			keyboard := keyboards.CreateSuccessKeyboard("global_settings")
			b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, errorText, &keyboard)
			return
		}
	}

	b.setUserState(message.From.ID, types.StateIdle)

	// Проверяем, что settings не nil перед использованием
	if settings != nil {
		logInfo("GlobalMarkup", fmt.Sprintf("Settings saved successfully. Returned value: %.2f%%", settings.GlobalMarkupPercent))
		text := fmt.Sprintf(MessageGlobalMarkupSet, settings.GlobalMarkupPercent)
		logInfo("GlobalMarkup", fmt.Sprintf("Showing confirmation with value: %.2f%%", settings.GlobalMarkupPercent))
		keyboard := keyboards.CreateSuccessKeyboard("global_settings")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
	} else {
		logError("GlobalMarkup", fmt.Errorf("settings is nil"))
		// Показываем сообщение с введенным значением, если нет данных от сервера
		text := fmt.Sprintf(MessageGlobalMarkupSetWithNote, markup)
		logInfo("GlobalMarkup", fmt.Sprintf("Showing confirmation with input value: %.2f%%", markup))
		keyboard := keyboards.CreateSuccessKeyboard("global_settings")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
	}
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
func (b *Bot) confirmCreateSubscription(userID, _ int64, userState *types.UserData) {
	log.Printf("User %d confirming subscription creation", userID)

	if userState.SubscriptionData == nil {
		log.Printf("ERROR: SubscriptionData is nil for user %d", userID)
		keyboard := keyboards.CreateSuccessKeyboard("manage_subscriptions")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, MessageDataNotFound, &keyboard)
		return
	}

	data := userState.SubscriptionData
	log.Printf("Subscription data for user %d: %+v", userID, data)

	// Создаем подписку через API
	req := api.CreateSubscriptionRequest{
		ServiceName:  data.ServiceName,
		BasePrice:    data.BasePrice,
		BaseCurrency: data.BaseCurrency,
		PeriodDays:   data.PeriodDays,
	}

	log.Printf("Creating subscription request: %+v", req)

	subscription, err := b.Context.APIClient.CreateSubscription(req)
	if err != nil {
		log.Printf("ERROR: Failed to create subscription for user %d: %v", userID, err)
		errorText := fmt.Sprintf(MessageSubscriptionCreateError, handleAPIError(err, "CreateSubscription"))
		keyboard := keyboards.CreateSuccessKeyboard("manage_subscriptions")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, errorText, &keyboard)
		return
	}

	log.Printf("Subscription created successfully: %+v", subscription)

	text := fmt.Sprintf(MessageSubscriptionCreated, subscription.ServiceName, subscription.BasePrice, subscription.BaseCurrency, subscription.PeriodDays, subscription.ID)

	// Сбрасываем состояние
	b.resetUserState(userID)

	keyboard := keyboards.CreateSuccessKeyboard("manage_subscriptions")
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// confirmCreateUser подтверждает создание пользователя
func (b *Bot) confirmCreateUser(userID, _ int64, userState *types.UserData) {
	log.Printf("User %d confirming user creation", userID)

	if userState.UserCreateData == nil {
		log.Printf("ERROR: UserCreateData is nil for user %d", userID)
		keyboard := keyboards.CreateSuccessKeyboard("manage_users")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, MessageDataNotFound, &keyboard)
		return
	}

	data := userState.UserCreateData
	log.Printf("User create data for user %d: %+v", userID, data)

	// Создаем пользователя через API
	req := api.CreateUserRequest{
		TGID:     data.TGID,
		Username: data.Username,
		Fullname: data.Fullname,
		IsAdmin:  false, // По умолчанию не админ
	}

	log.Printf("Creating user request: %+v", req)

	user, err := b.Context.APIClient.CreateUser(req)
	if err != nil {
		log.Printf("ERROR: Failed to create user for user %d: %v", userID, err)
		errorText := fmt.Sprintf(MessageUserCreateError, handleAPIError(err, "CreateUser"))
		keyboard := keyboards.CreateSuccessKeyboard("manage_users")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, errorText, &keyboard)
		return
	}

	log.Printf("User created successfully: %+v", user)

	usernameText := formatUsername(user.Username)

	text := fmt.Sprintf(MessageUserCreated, user.Fullname, user.TGID, usernameText, user.ID)

	// Сбрасываем состояние
	b.resetUserState(userID)

	keyboard := keyboards.CreateSuccessKeyboard("manage_users")
	b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
}

// cancelCurrentOperation отменяет текущую операцию
func (b *Bot) cancelCurrentOperation(userID, chatID int64) {
	userState := b.getUserState(userID)
	menuContext := userState.CurrentMenuContext

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.SubscriptionData = nil
	userState.UserCreateData = nil
	userState.CurrentEntityID = ""

	// Возвращаемся в соответствующее меню в зависимости от контекста
	if userState.CurrentMessageID != 0 {
		switch menuContext {
		case "subscriptions":
			b.showSubscriptionManagementEdit(userState.CurrentChatID, userState.CurrentMessageID)
		case "users":
			b.showUserManagementEdit(userState.CurrentChatID, userState.CurrentMessageID)
		case "global_settings":
			b.showGlobalSettingsEdit(userState.CurrentChatID, userState.CurrentMessageID)
		default:
			b.showMainMenuEdit(userState.CurrentChatID, userState.CurrentMessageID, userID)
		}
	} else {
		// Если нет сохраненного messageID, отправляем новое сообщение
		// (это не должно происходить в нормальном процессе)
		switch menuContext {
		case "subscriptions":
			b.showSubscriptionManagement(chatID)
		case "users":
			b.showUserManagement(chatID)
		case "global_settings":
			b.showGlobalSettings(chatID)
		default:
			b.showMainMenu(chatID, userID)
		}
	}
}

// handleStepBack обрабатывает возвращение на шаг назад
func (b *Bot) handleStepBack(userID, chatID int64) {
	userState := b.getUserState(userID)

	switch userState.State {
	// Создание подписки
	case types.StateAwaitingSubscriptionPrice:
		userState.State = types.StateAwaitingSubscriptionName
		text := "📝 Создание новой подписки\n\n**Шаг 1/4:** Введите название сервиса\n\n*Например: Netflix, Spotify*"
		keyboard := keyboards.CreateProcessKeyboard("start")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)

	case types.StateAwaitingSubscriptionCurrency:
		userState.State = types.StateAwaitingSubscriptionPrice
		text := fmt.Sprintf("📝 Создание новой подписки\n\n**Шаг 2/4:** Введите базовую стоимость\n\n✅ Название: %s\n\n*Например: 9.99*", userState.SubscriptionData.ServiceName)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)

	case types.StateAwaitingSubscriptionPeriod:
		userState.State = types.StateAwaitingSubscriptionCurrency
		text := fmt.Sprintf("📝 Создание новой подписки\n\n**Шаг 3/4:** Выберите валюту\n\n✅ Название: %s\n✅ Цена: %.2f", userState.SubscriptionData.ServiceName, userState.SubscriptionData.BasePrice)
		keyboard := keyboards.CurrencyKeyboardWithNav()
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)

	// Создание пользователя
	case types.StateAwaitingUserTGID:
		userState.State = types.StateAwaitingUserFullname
		text := "👤 Создание нового пользователя\n\n**Шаг 1/3:** Введите ФИО пользователя\n\n*Например: Иван Петров*"
		keyboard := keyboards.CreateProcessKeyboard("start")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)

	case types.StateAwaitingUserUsername:
		userState.State = types.StateAwaitingUserTGID
		text := fmt.Sprintf("👤 Создание нового пользователя\n\n**Шаг 2/3:** Введите Telegram ID пользователя\n\n✅ ФИО: %s\n\n*Например: 123456789*", userState.UserCreateData.Fullname)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)

	default:
		// Если не можем вернуться назад, отменяем операцию
		b.cancelCurrentOperation(userID, chatID)
	}
}
