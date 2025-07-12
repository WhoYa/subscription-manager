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

// Bot основная структура бота
type Bot struct {
	API     *tgbotapi.BotAPI
	Context *types.BotContext
}

// NewBot создает новый экземпляр бота
func NewBot(token, apiBaseURL string, adminUserIDs []int64) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	botAPI.Debug = false

	// Создаем API клиент
	apiClient := api.NewClient(apiBaseURL)

	context := &types.BotContext{
		Bot:          botAPI,
		APIClient:    apiClient,
		APIBaseURL:   apiBaseURL,
		UserStates:   make(map[int64]*types.UserData),
		AdminUserIDs: adminUserIDs,
	}

	return &Bot{
		API:     botAPI,
		Context: context,
	}, nil
}

// Start запускает бота
func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}

	return nil
}

// handleUpdate обрабатывает входящие обновления
func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		b.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update.CallbackQuery)
	}
}

// handleMessage обрабатывает входящие сообщения
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	if !b.isAdmin(message.From.ID) {
		b.sendMessage(message.Chat.ID, "❌ У вас нет прав доступа к этому боту.")
		return
	}

	// Получаем состояние пользователя
	userState := b.getUserState(message.From.ID)

	// Обрабатываем команды
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	// Обрабатываем сообщения в зависимости от состояния
	switch userState.State {
	case types.StateAwaitingSubscriptionName:
		b.handleSubscriptionNameInput(message)
	case types.StateAwaitingSubscriptionPrice:
		b.handleSubscriptionPriceInput(message)
	case types.StateAwaitingSubscriptionPeriod:
		b.handleSubscriptionPeriodInput(message)
	case types.StateAwaitingUserFullname:
		b.handleUserFullnameInput(message)
	case types.StateAwaitingUserTGID:
		b.handleUserTGIDInput(message)
	case types.StateAwaitingUserUsername:
		b.handleUserUsernameInput(message)
	case types.StateAwaitingGlobalMarkup:
		b.handleGlobalMarkupInput(message)
	default:
		b.sendMessage(message.Chat.ID, "Используйте /start для начала работы или /help для получения справки.")
	}
}

// handleCommand обрабатывает команды
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.handleStartCommand(message)
	case "help":
		b.handleHelpCommand(message)
	case "menu":
		b.showMainMenu(message.Chat.ID)
	default:
		b.sendMessage(message.Chat.ID, "Неизвестная команда. Используйте /help для получения справки.")
	}
}

// handleStartCommand обрабатывает команду /start
func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	welcomeText := fmt.Sprintf(`
👋 Добро пожаловать в бота управления подписками!

Привет, %s!

Этот бот позволяет администраторам управлять подписками и пользователями.

Доступные команды:
/menu - Главное меню
/help - Справка

Выберите действие из меню ниже:`, message.From.FirstName)

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ReplyMarkup = keyboards.MainAdminKeyboard()
	b.API.Send(msg)

	// Сбрасываем состояние пользователя
	b.setUserState(message.From.ID, types.StateIdle)
}

// handleHelpCommand обрабатывает команду /help
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `
📖 Справка по боту управления подписками

Основные функции:
📝 Управление подписками - создание и редактирование подписок
👥 Управление пользователями - создание профилей пользователей
⚙️ Глобальные настройки - настройка процента надбавки
📊 Аналитика - просмотр статистики прибыли

Команды:
/start - Главное меню
/menu - Показать главное меню
/help - Эта справка

Для начала работы используйте /start`

	b.sendMessage(message.Chat.ID, helpText)
}

// handleCallbackQuery обрабатывает callback запросы
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	if !b.isAdmin(query.From.ID) {
		b.answerCallbackQuery(query.ID, "❌ У вас нет прав доступа.")
		return
	}

	b.answerCallbackQuery(query.ID, "")

	switch query.Data {
	case "main_menu":
		b.showMainMenu(query.Message.Chat.ID)
	case "manage_subscriptions":
		b.showSubscriptionManagement(query.Message.Chat.ID)
	case "manage_users":
		b.showUserManagement(query.Message.Chat.ID)
	case "global_settings":
		b.showGlobalSettings(query.Message.Chat.ID)
	case "analytics":
		b.showAnalytics(query.Message.Chat.ID)
	case "create_subscription":
		b.startCreateSubscription(query.From.ID, query.Message.Chat.ID)
	case "create_user":
		b.startCreateUser(query.From.ID, query.Message.Chat.ID)
	case "edit_global_markup":
		b.startEditGlobalMarkup(query.From.ID, query.Message.Chat.ID)
	case "list_subscriptions":
		b.handleListSubscriptions(query.Message.Chat.ID)
	case "list_users":
		b.handleListUsers(query.Message.Chat.ID)
	case "analytics_total":
		b.handleAnalyticsTotal(query.Message.Chat.ID, query.From.ID)
	case "analytics_monthly":
		b.handleAnalyticsMonthly(query.Message.Chat.ID, query.From.ID)
	case "cancel":
		b.cancelCurrentOperation(query.From.ID, query.Message.Chat.ID)
	default:
		if strings.HasPrefix(query.Data, "currency_") {
			b.handleCurrencySelection(query)
		} else if strings.HasPrefix(query.Data, "confirm_") {
			b.handleConfirmation(query)
		} else {
			b.sendMessage(query.Message.Chat.ID, "Функция пока не реализована.")
		}
	}
}

// showMainMenu показывает главное меню
func (b *Bot) showMainMenu(chatID int64) {
	text := `
🏠 Главное меню администратора

Выберите нужный раздел:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.MainAdminKeyboard()
	b.API.Send(msg)
}

// showSubscriptionManagement показывает меню управления подписками
func (b *Bot) showSubscriptionManagement(chatID int64) {
	text := `
📝 Управление подписками

Выберите действие:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.SubscriptionManagementKeyboard()
	b.API.Send(msg)
}

// showUserManagement показывает меню управления пользователями
func (b *Bot) showUserManagement(chatID int64) {
	text := `
👥 Управление пользователями

Выберите действие:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.UserManagementKeyboard()
	b.API.Send(msg)
}

// showGlobalSettings показывает глобальные настройки
func (b *Bot) showGlobalSettings(chatID int64) {
	// Получаем текущие настройки из API
	settings, err := b.Context.APIClient.GetGlobalSettings()
	if err != nil {
		text := fmt.Sprintf(`
⚙️ Глобальные настройки

❌ Ошибка при загрузке настроек: %v

Выберите действие:`, err)

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📝 Изменить надбавку", "edit_global_markup"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "global_settings"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu"),
			),
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	text := fmt.Sprintf(`
⚙️ Глобальные настройки

Текущая глобальная надбавка: %.2f%%

Выберите действие:`, settings.GlobalMarkupPercent)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Изменить надбавку", "edit_global_markup"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "global_settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

// showAnalytics показывает меню аналитики
func (b *Bot) showAnalytics(chatID int64) {
	text := `
📊 Аналитика и статистика

Выберите тип отчета:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.AnalyticsKeyboard()
	b.API.Send(msg)
}

// Utility functions

// isAdmin проверяет, является ли пользователь администратором
func (b *Bot) isAdmin(userID int64) bool {
	for _, adminID := range b.Context.AdminUserIDs {
		if adminID == userID {
			return true
		}
	}
	return false
}

// getUserState получает состояние пользователя
func (b *Bot) getUserState(userID int64) *types.UserData {
	if state, exists := b.Context.UserStates[userID]; exists {
		return state
	}

	// Создаем новое состояние если его нет
	newState := &types.UserData{
		State: types.StateIdle,
	}
	b.Context.UserStates[userID] = newState
	return newState
}

// setUserState устанавливает состояние пользователя
func (b *Bot) setUserState(userID int64, state types.UserState) {
	userState := b.getUserState(userID)
	userState.State = state
}

// sendMessage отправляет текстовое сообщение
func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.API.Send(msg)
}

// answerCallbackQuery отвечает на callback запрос
func (b *Bot) answerCallbackQuery(queryID, text string) {
	callback := tgbotapi.NewCallback(queryID, text)
	b.API.Request(callback)
}
