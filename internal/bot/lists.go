package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/WhoYa/subscription-manager/internal/bot/keyboards"
)

// Subscriptions list handlers

// handleListSubscriptions показывает список подписок
func (b *Bot) handleListSubscriptions(chatID int64) {
	subscriptions, err := b.Context.APIClient.GetSubscriptions(25, 0)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке подписок: %v", err))
		return
	}

	if len(subscriptions) == 0 {
		text := `
📋 Список подписок

📭 Подписки не найдены.

Создайте первую подписку!`

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Создать подписку", "create_subscription"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_subscriptions"),
			),
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	var textBuilder strings.Builder
	textBuilder.WriteString("📋 Список подписок\n\n")

	for i, sub := range subscriptions {
		status := "✅"
		if !sub.IsActive {
			status = "❌"
		}

		textBuilder.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, status, sub.ServiceName))
		textBuilder.WriteString(fmt.Sprintf("   💰 %.2f %s, период: %d дней\n",
			sub.BasePrice, sub.BaseCurrency, sub.PeriodDays))
		textBuilder.WriteString(fmt.Sprintf("   🆔 ID: %s\n\n", sub.ID))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать подписку", "create_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Редактировать", "edit_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "list_subscriptions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_subscriptions"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, textBuilder.String())
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

// Users list handlers

// handleListUsers показывает список пользователей
func (b *Bot) handleListUsers(chatID int64) {
	users, err := b.Context.APIClient.GetUsers(25, 0)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке пользователей: %v", err))
		return
	}

	if len(users) == 0 {
		text := `
📋 Список пользователей

📭 Пользователи не найдены.

Создайте первого пользователя!`

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Создать пользователя", "create_user"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_users"),
			),
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	var textBuilder strings.Builder
	textBuilder.WriteString("📋 Список пользователей\n\n")

	for i, user := range users {
		adminStatus := ""
		if user.IsAdmin {
			adminStatus = " 👑"
		}

		usernameText := ""
		if user.Username != "" {
			usernameText = fmt.Sprintf(" (@%s)", user.Username)
		}

		textBuilder.WriteString(fmt.Sprintf("%d. %s%s%s\n", i+1, user.Fullname, usernameText, adminStatus))
		textBuilder.WriteString(fmt.Sprintf("   🆔 TG ID: %d\n", user.TGID))
		textBuilder.WriteString(fmt.Sprintf("   🆔 ID: %s\n\n", user.ID))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать пользователя", "create_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Редактировать", "edit_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "list_users"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_users"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, textBuilder.String())
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

// Analytics handlers

// handleAnalyticsTotal показывает общую статистику
func (b *Bot) handleAnalyticsTotal(chatID, adminUserID int64) {
	adminUser, err := b.getAdminUser(adminUserID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ %v", err))
		return
	}

	stats, err := b.Context.APIClient.GetTotalProfit(adminUser.ID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке статистики: %v", err))
		return
	}

	text := fmt.Sprintf(`
📈 Общая статистика прибыли

💰 Общая прибыль: %.2f руб.
🧾 Количество платежей: %d
📊 Средняя прибыль с платежа: %.2f руб.`,
		stats.TotalProfit, stats.PaymentCount, stats.AverageProfit)

	keyboard := keyboards.BackKeyboard("analytics")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

// handleAnalyticsMonthly показывает статистику за текущий месяц
func (b *Bot) handleAnalyticsMonthly(chatID, adminUserID int64) {
	adminUser, err := b.getAdminUser(adminUserID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ %v", err))
		return
	}

	// Получаем текущую дату для показа статистики за текущий месяц
	// Для примера возьмем 2024 год, 7 месяц
	stats, err := b.Context.APIClient.GetMonthlyProfit(adminUser.ID, 2024, 7)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке статистики: %v", err))
		return
	}

	text := fmt.Sprintf(`
📅 Статистика за июль 2024

💰 Прибыль за месяц: %.2f руб.
🧾 Количество платежей: %d
📊 Средняя прибыль с платежа: %.2f руб.`,
		stats.TotalProfit, stats.PaymentCount, stats.AverageProfit)

	keyboard := keyboards.BackKeyboard("analytics")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

// handleAnalyticsUsers показывает аналитику по пользователям
func (b *Bot) handleAnalyticsUsers(chatID, adminUserID int64) {
	adminUser, err := b.getAdminUser(adminUserID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ %v", err))
		return
	}

	// Получаем статистику за последний месяц
	from, to := b.getLastMonthRange()
	stats, err := b.Context.APIClient.GetUserProfitStats(adminUser.ID, from, to)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке статистики: %v", err))
		return
	}

	if len(stats) == 0 {
		text := `👥 Аналитика по пользователям

📭 Нет данных за последний месяц.`
		keyboard := keyboards.BackKeyboard("analytics")
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	var textBuilder strings.Builder
	textBuilder.WriteString("👥 Аналитика по пользователям\n")
	textBuilder.WriteString(fmt.Sprintf("📅 Период: %s - %s\n\n", from[:10], to[:10]))

	totalProfit := 0.0
	totalPayments := 0

	for i, stat := range stats {
		textBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, stat.UserFullname))
		textBuilder.WriteString(fmt.Sprintf("   💰 Прибыль: %.2f руб.\n", stat.TotalProfit))
		textBuilder.WriteString(fmt.Sprintf("   🧾 Платежей: %d\n\n", stat.PaymentCount))

		totalProfit += stat.TotalProfit
		totalPayments += stat.PaymentCount
	}

	textBuilder.WriteString(fmt.Sprintf("📊 Итого: %.2f руб. (%d платежей)", totalProfit, totalPayments))

	keyboard := keyboards.BackKeyboard("analytics")
	msg := tgbotapi.NewMessage(chatID, textBuilder.String())
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

// handleAnalyticsSubscriptions показывает аналитику по подпискам
func (b *Bot) handleAnalyticsSubscriptions(chatID, adminUserID int64) {
	adminUser, err := b.getAdminUser(adminUserID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ %v", err))
		return
	}

	// Получаем статистику за последний месяц
	from, to := b.getLastMonthRange()
	stats, err := b.Context.APIClient.GetSubscriptionProfitStats(adminUser.ID, from, to)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при загрузке статистики: %v", err))
		return
	}

	if len(stats) == 0 {
		text := `📱 Аналитика по подпискам

📭 Нет данных за последний месяц.`
		keyboard := keyboards.BackKeyboard("analytics")
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	var textBuilder strings.Builder
	textBuilder.WriteString("📱 Аналитика по подпискам\n")
	textBuilder.WriteString(fmt.Sprintf("📅 Период: %s - %s\n\n", from[:10], to[:10]))

	totalProfit := 0.0
	totalPayments := 0

	for i, stat := range stats {
		textBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, stat.SubscriptionName))
		textBuilder.WriteString(fmt.Sprintf("   💰 Прибыль: %.2f руб.\n", stat.TotalProfit))
		textBuilder.WriteString(fmt.Sprintf("   🧾 Платежей: %d\n\n", stat.PaymentCount))

		totalProfit += stat.TotalProfit
		totalPayments += stat.PaymentCount
	}

	textBuilder.WriteString(fmt.Sprintf("📊 Итого: %.2f руб. (%d платежей)", totalProfit, totalPayments))

	keyboard := keyboards.BackKeyboard("analytics")
	msg := tgbotapi.NewMessage(chatID, textBuilder.String())
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

// getLastMonthRange возвращает диапазон дат за последний месяц в формате RFC3339
func (b *Bot) getLastMonthRange() (string, string) {
	now := time.Now()
	from := now.AddDate(0, -1, 0) // месяц назад
	to := now

	return from.Format(time.RFC3339), to.Format(time.RFC3339)
}

// Edit versions of list functions

// handleListSubscriptionsEdit показывает список подписок через редактирование сообщения
func (b *Bot) handleListSubscriptionsEdit(chatID int64, messageID int) {
	subscriptions, err := b.Context.APIClient.GetSubscriptions(25, 0)
	if err != nil {
		b.editMessage(chatID, messageID, fmt.Sprintf("❌ Ошибка при загрузке подписок: %v", err), nil)
		return
	}

	if len(subscriptions) == 0 {
		text := `📋 Список подписок

📭 Подписки не найдены.

Создайте первую подписку!`

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Создать подписку", "create_subscription"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_subscriptions_edit"),
			),
		)

		b.editMessage(chatID, messageID, text, &keyboard)
		return
	}

	var textBuilder strings.Builder
	textBuilder.WriteString("📋 Список подписок\n\n")

	for i, sub := range subscriptions {
		status := "✅"
		if !sub.IsActive {
			status = "❌"
		}

		textBuilder.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, status, sub.ServiceName))
		textBuilder.WriteString(fmt.Sprintf("   💰 %.2f %s, период: %d дней\n",
			sub.BasePrice, sub.BaseCurrency, sub.PeriodDays))
		textBuilder.WriteString(fmt.Sprintf("   🆔 ID: %s\n\n", sub.ID))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать подписку", "create_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Редактировать", "edit_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "list_subscriptions_edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_subscriptions_edit"),
		),
	)

	b.editMessage(chatID, messageID, textBuilder.String(), &keyboard)
}

// handleListUsersEdit показывает список пользователей через редактирование сообщения
func (b *Bot) handleListUsersEdit(chatID int64, messageID int) {
	users, err := b.Context.APIClient.GetUsers(25, 0)
	if err != nil {
		b.editMessage(chatID, messageID, fmt.Sprintf("❌ Ошибка при загрузке пользователей: %v", err), nil)
		return
	}

	if len(users) == 0 {
		text := `📋 Список пользователей

📭 Пользователи не найдены.

Создайте первого пользователя!`

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Создать пользователя", "create_user"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_users_edit"),
			),
		)

		b.editMessage(chatID, messageID, text, &keyboard)
		return
	}

	var textBuilder strings.Builder
	textBuilder.WriteString("📋 Список пользователей\n\n")

	for i, user := range users {
		adminStatus := ""
		if user.IsAdmin {
			adminStatus = " 👑"
		}

		usernameText := ""
		if user.Username != "" {
			usernameText = fmt.Sprintf(" (@%s)", user.Username)
		}

		textBuilder.WriteString(fmt.Sprintf("%d. %s%s%s\n", i+1, user.Fullname, usernameText, adminStatus))
		textBuilder.WriteString(fmt.Sprintf("   🆔 TG ID: %d\n", user.TGID))
		textBuilder.WriteString(fmt.Sprintf("   🆔 ID: %s\n\n", user.ID))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать пользователя", "create_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Редактировать", "edit_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "list_users_edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "manage_users_edit"),
		),
	)

	b.editMessage(chatID, messageID, textBuilder.String(), &keyboard)
}
