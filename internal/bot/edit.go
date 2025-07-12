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

// handleEditSubscription показывает список подписок для редактирования
func (b *Bot) handleEditSubscription(chatID int64, messageID int) {
	subscriptions, err := b.Context.APIClient.GetSubscriptions(25, 0)
	if err != nil {
		b.sendErrorMessage(chatID, messageID, err, "manage_subscriptions")
		return
	}

	if len(subscriptions) == 0 {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(ButtonCreateSubscription, "create_subscription"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(ButtonBack, "manage_subscriptions"),
			),
		)

		b.editMessage(chatID, messageID, MessageEditSubscriptionEmpty, &keyboard)
		return
	}

	// Создаем клавиатуру с подписками
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, sub := range subscriptions {
		status := getSubscriptionStatus(sub.IsActive)
		buttonText := fmt.Sprintf("%s %s (%.2f %s)", status, sub.ServiceName, sub.BasePrice, sub.BaseCurrency)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, fmt.Sprintf("edit_sub_%s", sub.ID)),
		))
	}

	// Добавляем кнопку "Назад"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ButtonBack, "manage_subscriptions"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	b.editMessage(chatID, messageID, MessageEditSubscriptionTitle, &keyboard)
}

// handleEditUser показывает список пользователей для редактирования
func (b *Bot) handleEditUser(chatID int64, messageID int) {
	users, err := b.Context.APIClient.GetUsers(25, 0)
	if err != nil {
		b.sendErrorMessage(chatID, messageID, err, "manage_users")
		return
	}

	if len(users) == 0 {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(ButtonCreateUser, "create_user"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(ButtonBack, "manage_users"),
			),
		)

		b.editMessage(chatID, messageID, MessageEditUserEmpty, &keyboard)
		return
	}

	// Создаем клавиатуру с пользователями
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, user := range users {
		adminStatus := ""
		if user.IsAdmin {
			adminStatus = " 👑"
		}

		buttonText := fmt.Sprintf("%s%s", user.Fullname, adminStatus)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, fmt.Sprintf("edit_user_%s", user.ID)),
		))
	}

	// Добавляем кнопку "Назад"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ButtonBack, "manage_users"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	b.editMessage(chatID, messageID, MessageEditUserTitle, &keyboard)
}

// showSubscriptionEditMenu показывает меню редактирования конкретной подписки
func (b *Bot) showSubscriptionEditMenu(chatID int64, messageID int, subscriptionID string) {
	subscription, err := b.Context.APIClient.GetSubscription(subscriptionID)
	if err != nil {
		b.sendErrorMessage(chatID, messageID, err, "edit_subscription")
		return
	}

	status := getSubscriptionStatus(subscription.IsActive)

	text := fmt.Sprintf(MessageEditSubscriptionMenu, subscription.ServiceName, subscription.BasePrice, subscription.BaseCurrency, subscription.PeriodDays, status)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonEditName, fmt.Sprintf("edit_sub_name_%s", subscriptionID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonEditPrice, fmt.Sprintf("edit_sub_price_%s", subscriptionID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonEditCurrency, fmt.Sprintf("edit_sub_currency_%s", subscriptionID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonEditPeriod, fmt.Sprintf("edit_sub_period_%s", subscriptionID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonToggleStatus, fmt.Sprintf("toggle_sub_status_%s", subscriptionID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonBack, "edit_subscription"),
		),
	)

	b.editMessage(chatID, messageID, text, &keyboard)
}

// showUserEditMenu показывает меню редактирования конкретного пользователя
func (b *Bot) showUserEditMenu(chatID int64, messageID int, userID string) {
	user, err := b.Context.APIClient.GetUser(userID)
	if err != nil {
		b.sendErrorMessage(chatID, messageID, err, "edit_user")
		return
	}

	adminStatus := getUserRoleStatus(user.IsAdmin)
	usernameText := formatUsername(user.Username)

	text := fmt.Sprintf(MessageEditUserMenu, user.Fullname, user.TGID, usernameText, adminStatus)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonEditFullname, fmt.Sprintf("edit_user_fullname_%s", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonEditUsername, fmt.Sprintf("edit_user_username_%s", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonToggleRole, fmt.Sprintf("toggle_user_admin_%s", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonBack, "edit_user"),
		),
	)

	b.editMessage(chatID, messageID, text, &keyboard)
}

// Обработчики ввода для редактирования подписок
func (b *Bot) handleEditSubscriptionNameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)
	if userState.EditData == nil {
		b.sendMessage(message.Chat.ID, "❌ Ошибка: данные для редактирования не найдены.")
		return
	}

	newName := strings.TrimSpace(message.Text)
	if newName == "" {
		b.sendMessage(message.Chat.ID, "❌ Название не может быть пустым. Попробуйте снова:")
		return
	}

	// Обновляем подписку
	req := api.UpdateSubscriptionRequest{
		ServiceName: &newName,
	}

	subscription, err := b.Context.APIClient.UpdateSubscription(userState.EditData.EntityID, req)
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Ошибка при обновлении подписки: %v", err))
		return
	}

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.EditData = nil

	text := fmt.Sprintf("✅ Название подписки обновлено: %s", subscription.ServiceName)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard(fmt.Sprintf("edit_sub_%s", subscription.ID))
	b.API.Send(msg)
}

func (b *Bot) handleEditSubscriptionPriceInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)
	if userState.EditData == nil {
		b.sendMessage(message.Chat.ID, "❌ Ошибка: данные для редактирования не найдены.")
		return
	}

	priceStr := strings.TrimSpace(message.Text)
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		b.sendMessage(message.Chat.ID, "❌ Неверный формат цены. Введите число больше 0:")
		return
	}

	// Обновляем подписку
	req := api.UpdateSubscriptionRequest{
		BasePrice: &price,
	}

	subscription, err := b.Context.APIClient.UpdateSubscription(userState.EditData.EntityID, req)
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Ошибка при обновлении подписки: %v", err))
		return
	}

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.EditData = nil

	text := fmt.Sprintf("✅ Цена подписки обновлена: %.2f %s", subscription.BasePrice, subscription.BaseCurrency)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard(fmt.Sprintf("edit_sub_%s", subscription.ID))
	b.API.Send(msg)
}

func (b *Bot) handleEditSubscriptionPeriodInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)
	if userState.EditData == nil {
		b.sendMessage(message.Chat.ID, "❌ Ошибка: данные для редактирования не найдены.")
		return
	}

	periodStr := strings.TrimSpace(message.Text)
	period, err := strconv.Atoi(periodStr)
	if err != nil || period <= 0 {
		b.sendMessage(message.Chat.ID, "❌ Неверный формат периода. Введите целое число больше 0:")
		return
	}

	// Обновляем подписку
	req := api.UpdateSubscriptionRequest{
		PeriodDays: &period,
	}

	subscription, err := b.Context.APIClient.UpdateSubscription(userState.EditData.EntityID, req)
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Ошибка при обновлении подписки: %v", err))
		return
	}

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.EditData = nil

	text := fmt.Sprintf("✅ Период подписки обновлен: %d дней", subscription.PeriodDays)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard(fmt.Sprintf("edit_sub_%s", subscription.ID))
	b.API.Send(msg)
}

// Обработчики ввода для редактирования пользователей
func (b *Bot) handleEditUserFullnameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)
	if userState.EditData == nil {
		b.sendMessage(message.Chat.ID, "❌ Ошибка: данные для редактирования не найдены.")
		return
	}

	newFullname := strings.TrimSpace(message.Text)
	if newFullname == "" {
		b.sendMessage(message.Chat.ID, "❌ ФИО не может быть пустым. Попробуйте снова:")
		return
	}

	// Обновляем пользователя
	req := api.UpdateUserRequest{
		Fullname: &newFullname,
	}

	user, err := b.Context.APIClient.UpdateUser(userState.EditData.EntityID, req)
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Ошибка при обновлении пользователя: %v", err))
		return
	}

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.EditData = nil

	text := fmt.Sprintf("✅ ФИО пользователя обновлено: %s", user.Fullname)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard(fmt.Sprintf("edit_user_%s", user.ID))
	b.API.Send(msg)
}

func (b *Bot) handleEditUserUsernameInput(message *tgbotapi.Message) {
	userState := b.getUserState(message.From.ID)
	if userState.EditData == nil {
		b.sendMessage(message.Chat.ID, "❌ Ошибка: данные для редактирования не найдены.")
		return
	}

	newUsername := strings.TrimSpace(message.Text)
	// Username может быть пустым

	// Обновляем пользователя
	req := api.UpdateUserRequest{
		Username: &newUsername,
	}

	user, err := b.Context.APIClient.UpdateUser(userState.EditData.EntityID, req)
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Ошибка при обновлении пользователя: %v", err))
		return
	}

	// Сбрасываем состояние
	userState.State = types.StateIdle
	userState.EditData = nil

	usernameText := "не указан"
	if user.Username != "" {
		usernameText = "@" + user.Username
	}

	text := fmt.Sprintf("✅ Username пользователя обновлен: %s", usernameText)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard(fmt.Sprintf("edit_user_%s", user.ID))
	b.API.Send(msg)
}

// handleEditSubscriptionCallback обрабатывает callback'ы редактирования подписок
func (b *Bot) handleEditSubscriptionCallback(query *tgbotapi.CallbackQuery) {
	parts := strings.Split(query.Data, "_")
	if len(parts) < 3 {
		b.sendMessage(query.Message.Chat.ID, "❌ Неверный формат callback.")
		return
	}

	action := parts[2]
	subscriptionID := parts[len(parts)-1]

	switch action {
	case "name":
		b.startEditSubscriptionName(query.From.ID, query.Message.Chat.ID, subscriptionID)
	case "price":
		b.startEditSubscriptionPrice(query.From.ID, query.Message.Chat.ID, subscriptionID)
	case "currency":
		b.startEditSubscriptionCurrency(query.From.ID, query.Message.Chat.ID, subscriptionID)
	case "period":
		b.startEditSubscriptionPeriod(query.From.ID, query.Message.Chat.ID, subscriptionID)
	default:
		// Показываем меню редактирования подписки
		b.showSubscriptionEditMenu(query.Message.Chat.ID, query.Message.MessageID, subscriptionID)
	}
}

// handleEditUserCallback обрабатывает callback'ы редактирования пользователей
func (b *Bot) handleEditUserCallback(query *tgbotapi.CallbackQuery) {
	parts := strings.Split(query.Data, "_")
	if len(parts) < 3 {
		b.sendMessage(query.Message.Chat.ID, "❌ Неверный формат callback.")
		return
	}

	action := parts[2]
	userID := parts[len(parts)-1]

	switch action {
	case "fullname":
		b.startEditUserFullname(query.From.ID, query.Message.Chat.ID, userID)
	case "username":
		b.startEditUserUsername(query.From.ID, query.Message.Chat.ID, userID)
	default:
		// Показываем меню редактирования пользователя
		b.showUserEditMenu(query.Message.Chat.ID, query.Message.MessageID, userID)
	}
}

// handleToggleCallback обрабатывает переключение статусов
func (b *Bot) handleToggleCallback(query *tgbotapi.CallbackQuery) {
	parts := strings.Split(query.Data, "_")
	if len(parts) < 4 {
		b.sendMessage(query.Message.Chat.ID, "❌ Неверный формат callback.")
		return
	}

	entityType := parts[1]
	action := parts[2]
	entityID := parts[len(parts)-1]

	switch entityType {
	case "sub":
		if action == "status" {
			b.toggleSubscriptionStatus(query.Message.Chat.ID, entityID)
		}
	case "user":
		if action == "admin" {
			b.toggleUserAdminStatus(query.Message.Chat.ID, entityID)
		}
	default:
		b.sendMessage(query.Message.Chat.ID, "❌ Неизвестный тип сущности.")
	}
}

// Методы для начала редактирования подписок
func (b *Bot) startEditSubscriptionName(userID, chatID int64, subscriptionID string) {
	userState := b.getUserState(userID)
	userState.State = types.StateEditingSubscriptionName
	userState.EditData = &types.EditData{
		EntityType: "subscription",
		EntityID:   subscriptionID,
	}

	b.sendSimpleMessage(chatID, MessageEditSubscriptionNamePrompt)
}

func (b *Bot) startEditSubscriptionPrice(userID, chatID int64, subscriptionID string) {
	userState := b.getUserState(userID)
	userState.State = types.StateEditingSubscriptionPrice
	userState.EditData = &types.EditData{
		EntityType: "subscription",
		EntityID:   subscriptionID,
	}

	b.sendSimpleMessage(chatID, MessageEditSubscriptionPricePrompt)
}

func (b *Bot) startEditSubscriptionCurrency(userID, chatID int64, subscriptionID string) {
	userState := b.getUserState(userID)
	userState.State = types.StateEditingSubscriptionCurrency
	userState.EditData = &types.EditData{
		EntityType: "subscription",
		EntityID:   subscriptionID,
	}

	text := "💱 Выберите новую валюту:"
	keyboard := keyboards.CurrencyKeyboard()
	b.sendMessageWithKeyboard(chatID, text, &keyboard)
}

func (b *Bot) startEditSubscriptionPeriod(userID, chatID int64, subscriptionID string) {
	userState := b.getUserState(userID)
	userState.State = types.StateEditingSubscriptionPeriod
	userState.EditData = &types.EditData{
		EntityType: "subscription",
		EntityID:   subscriptionID,
	}

	b.sendSimpleMessage(chatID, MessageEditSubscriptionPeriodPrompt)
}

// Методы для начала редактирования пользователей
func (b *Bot) startEditUserFullname(userID, chatID int64, targetUserID string) {
	userState := b.getUserState(userID)
	userState.State = types.StateEditingUserFullname
	userState.EditData = &types.EditData{
		EntityType: "user",
		EntityID:   targetUserID,
	}

	b.sendSimpleMessage(chatID, MessageEditUserFullnamePrompt)
}

func (b *Bot) startEditUserUsername(userID, chatID int64, targetUserID string) {
	userState := b.getUserState(userID)
	userState.State = types.StateEditingUserUsername
	userState.EditData = &types.EditData{
		EntityType: "user",
		EntityID:   targetUserID,
	}

	b.sendSimpleMessage(chatID, MessageEditUserUsernamePrompt)
}

// Методы для переключения статусов
func (b *Bot) toggleSubscriptionStatus(chatID int64, subscriptionID string) {
	subscription, err := b.Context.APIClient.GetSubscription(subscriptionID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке подписки: %v", err))
		return
	}

	newStatus := !subscription.IsActive
	req := api.UpdateSubscriptionRequest{
		IsActive: &newStatus,
	}

	updatedSubscription, err := b.Context.APIClient.UpdateSubscription(subscriptionID, req)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при обновлении статуса: %v", err))
		return
	}

	statusText := "активирована"
	if !updatedSubscription.IsActive {
		statusText = "деактивирована"
	}

	text := fmt.Sprintf("✅ Подписка %s %s", updatedSubscription.ServiceName, statusText)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard(fmt.Sprintf("edit_sub_%s", subscriptionID))
	b.API.Send(msg)
}

func (b *Bot) toggleUserAdminStatus(chatID int64, userID string) {
	user, err := b.Context.APIClient.GetUser(userID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке пользователя: %v", err))
		return
	}

	newStatus := !user.IsAdmin
	req := api.UpdateUserRequest{
		IsAdmin: &newStatus,
	}

	updatedUser, err := b.Context.APIClient.UpdateUser(userID, req)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при обновлении роли: %v", err))
		return
	}

	roleText := "администратор"
	if !updatedUser.IsAdmin {
		roleText = "обычный пользователь"
	}

	text := fmt.Sprintf("✅ Пользователь %s теперь %s", updatedUser.Fullname, roleText)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.BackKeyboard(fmt.Sprintf("edit_user_%s", userID))
	b.API.Send(msg)
}
