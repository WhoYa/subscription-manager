package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MainAdminKeyboard основное админ меню
func MainAdminKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Управление подписками", "manage_subscriptions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Управление пользователями", "manage_users"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Глобальные настройки", "global_settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Аналитика", "analytics"),
		),
	)
}

// SubscriptionManagementKeyboard меню управления подписками
func SubscriptionManagementKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать подписку", "create_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Редактировать подписку", "edit_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Список подписок", "list_subscriptions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu"),
		),
	)
}

// UserManagementKeyboard меню управления пользователями
func UserManagementKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать пользователя", "create_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Редактировать пользователя", "edit_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Список пользователей", "list_users"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu"),
		),
	)
}

// UserSubscriptionKeyboard меню управления подписками пользователя
func UserSubscriptionKeyboard(userID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Привязать подписку", fmt.Sprintf("user_add_sub_%s", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Настроить подписку", fmt.Sprintf("user_config_sub_%s", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖ Отписать", fmt.Sprintf("user_remove_sub_%s", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "edit_user"),
		),
	)
}

// AnalyticsKeyboard меню аналитики
func AnalyticsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📈 Общая прибыль", "analytics_total"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Прибыль за месяц", "analytics_monthly"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 По пользователям", "analytics_users"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 По подпискам", "analytics_subscriptions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu"),
		),
	)
}

// CurrencyKeyboard выбор валюты
func CurrencyKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💵 USD", "currency_USD"),
			tgbotapi.NewInlineKeyboardButtonData("💶 EUR", "currency_EUR"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "cancel"),
		),
	)
}

// ConfirmKeyboard подтверждение действия
func ConfirmKeyboard(action string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да", fmt.Sprintf("confirm_%s", action)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет", "cancel"),
		),
	)
}

// BackKeyboard простая кнопка назад
func BackKeyboard(backTo string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", backTo),
		),
	)
}

// PaginationKeyboard навигация по страницам
func PaginationKeyboard(action string, currentPage, totalPages int) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton

	if currentPage > 1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️", fmt.Sprintf("%s_page_%d", action, currentPage-1)))
	}

	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
		fmt.Sprintf("%d/%d", currentPage, totalPages),
		"noop",
	))

	if currentPage < totalPages {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("%s_page_%d", action, currentPage+1)))
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu"),
		),
	)
}
