package bot

import (
	"fmt"
	"strings"

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
	// Ищем админ пользователя в системе по его Telegram ID
	adminUser, err := b.Context.APIClient.FindUserByTGID(adminUserID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при поиске админ пользователя: %v", err))
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
	// Ищем админ пользователя в системе по его Telegram ID
	adminUser, err := b.Context.APIClient.FindUserByTGID(adminUserID)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при поиске админ пользователя: %v", err))
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
