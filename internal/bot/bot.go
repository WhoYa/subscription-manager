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

	log.Printf("Bot initialized with %d admin user(s): %v", len(adminUserIDs), adminUserIDs)

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
		b.sendSimpleMessage(message.Chat.ID, MessageNoAccess)
		return
	}

	// Получаем состояние пользователя
	userState := b.getUserState(message.From.ID)

	// Логируем входящее сообщение и состояние пользователя
	log.Printf("Received message from user %d (state: %s): %s", message.From.ID, userState.State, message.Text)

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
	case types.StateAwaitingSubscriptionCurrency:
		// Валюта выбирается через callback, но если пользователь отправил сообщение,
		// показываем подсказку о том, что нужно использовать кнопки
		log.Printf("User %d sent message in currency selection state: %s", message.From.ID, message.Text)
		b.sendSimpleMessage(message.Chat.ID, "💱 Пожалуйста, выберите валюту, используя кнопки ниже, или нажмите 'Отмена' для выхода.")
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
	case types.StateEditingSubscriptionName:
		b.handleEditSubscriptionNameInput(message)
	case types.StateEditingSubscriptionPrice:
		b.handleEditSubscriptionPriceInput(message)
	case types.StateEditingSubscriptionCurrency:
		// Валюта выбирается через callback, но если пользователь отправил сообщение,
		// показываем подсказку о том, что нужно использовать кнопки
		log.Printf("User %d sent message in currency editing state: %s", message.From.ID, message.Text)
		b.sendSimpleMessage(message.Chat.ID, "💱 Пожалуйста, выберите валюту, используя кнопки ниже, или нажмите 'Отмена' для выхода.")
	case types.StateEditingSubscriptionPeriod:
		b.handleEditSubscriptionPeriodInput(message)
	case types.StateEditingUserFullname:
		b.handleEditUserFullnameInput(message)
	case types.StateEditingUserUsername:
		b.handleEditUserUsernameInput(message)
	default:
		log.Printf("User %d sent message in unhandled state %s: %s", message.From.ID, userState.State, message.Text)
		b.sendSimpleMessage(message.Chat.ID, MessageUseStart)
	}
}

// handleCommand обрабатывает команды
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start", "menu":
		// При первом запуске создаем админа если его нет
		if message.Command() == "start" {
			// Проверяем, является ли пользователь админом в списке
			if b.isAdmin(message.From.ID) {
				// Проверяем, есть ли пользователь в БД
				user, err := b.Context.APIClient.FindUserByTGID(message.From.ID)
				if err != nil {
					// Пользователь не найден в БД, создаем его
					log.Printf("Admin user not found in database, creating...")
					_, err := b.getOrCreateAdminUser(message.From.ID, message.From.FirstName, message.From.LastName, message.From.UserName)
					if err != nil {
						log.Printf("Failed to create/find admin user: %v", err)
						b.sendSimpleMessage(message.Chat.ID, "❌ Ошибка при создании админского аккаунта")
						return
					}
				} else if !user.IsAdmin {
					// Пользователь найден, но не админ, повышаем его
					log.Printf("User exists but not admin, promoting to admin...")
					_, err := b.getOrCreateAdminUser(message.From.ID, message.From.FirstName, message.From.LastName, message.From.UserName)
					if err != nil {
						log.Printf("Failed to promote user to admin: %v", err)
						b.sendSimpleMessage(message.Chat.ID, "❌ Ошибка при обновлении прав доступа")
						return
					}
				}
				// Если пользователь уже админ, не делаем ничего
			}
		}
		b.showMainMenu(message.Chat.ID, message.From.ID)
	case "help":
		b.handleHelpCommand(message)
	default:
		b.sendSimpleMessage(message.Chat.ID, MessageUnknownCommand)
	}
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

	b.sendSimpleMessage(message.Chat.ID, helpText)
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
		b.showMainMenu(query.Message.Chat.ID, query.From.ID)
	case "main_menu_edit":
		b.showMainMenuEdit(query.Message.Chat.ID, query.Message.MessageID, query.From.ID)
	case "manage_subscriptions":
		b.showSubscriptionManagementEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "manage_subscriptions_edit":
		b.showSubscriptionManagementEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "manage_users":
		b.showUserManagementEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "manage_users_edit":
		b.showUserManagementEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "global_settings":
		b.showGlobalSettingsEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "global_settings_edit":
		b.showGlobalSettingsEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "analytics":
		b.showAnalyticsEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "analytics_edit":
		b.showAnalyticsEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "create_subscription":
		b.startCreateSubscription(query.From.ID, query.Message.Chat.ID, query.Message.MessageID)
	case "create_user":
		b.startCreateUser(query.From.ID, query.Message.Chat.ID, query.Message.MessageID)
	case "edit_global_markup":
		b.startEditGlobalMarkup(query.From.ID, query.Message.Chat.ID, query.Message.MessageID)
	case "list_subscriptions":
		b.handleListSubscriptions(query.Message.Chat.ID)
	case "list_subscriptions_edit":
		b.handleListSubscriptionsEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "list_users":
		b.handleListUsers(query.Message.Chat.ID)
	case "list_users_edit":
		b.handleListUsersEdit(query.Message.Chat.ID, query.Message.MessageID)
	case "analytics_total":
		b.handleAnalyticsTotal(query.Message.Chat.ID, query.From.ID)
	case "analytics_monthly":
		b.handleAnalyticsMonthly(query.Message.Chat.ID, query.From.ID)
	case "analytics_users":
		b.handleAnalyticsUsers(query.Message.Chat.ID, query.From.ID)
	case "analytics_subscriptions":
		b.handleAnalyticsSubscriptions(query.Message.Chat.ID, query.From.ID)
	case "edit_subscription":
		b.handleEditSubscription(query.Message.Chat.ID, query.Message.MessageID)
	case "edit_user":
		b.handleEditUser(query.Message.Chat.ID, query.Message.MessageID)
	case "cancel":
		b.cancelCurrentOperation(query.From.ID, query.Message.Chat.ID)
	case "step_back":
		b.handleStepBack(query.From.ID, query.Message.Chat.ID)
	default:
		if strings.HasPrefix(query.Data, "currency_") {
			b.handleCurrencySelection(query)
		} else if strings.HasPrefix(query.Data, "confirm_") {
			b.handleConfirmation(query)
		} else if strings.HasPrefix(query.Data, "edit_sub_") {
			b.handleEditSubscriptionCallback(query)
		} else if strings.HasPrefix(query.Data, "edit_user_") {
			b.handleEditUserCallback(query)
		} else if strings.HasPrefix(query.Data, "toggle_") {
			b.handleToggleCallback(query)
		} else {
			b.sendSimpleMessage(query.Message.Chat.ID, "Функция пока не реализована.")
		}
	}
}

// showMainMenu показывает главное меню
func (b *Bot) showMainMenu(chatID int64, userID ...int64) {
	var firstName string

	// Получаем имя пользователя из Telegram API
	if len(userID) > 0 {
		user, err := b.API.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: userID[0]}})
		if err == nil {
			firstName = user.FirstName
		}
	}

	// Если не удалось получить имя или не передан userID, используем по умолчанию
	if firstName == "" {
		firstName = "Администратор"
	}

	// Всегда показываем приветствие с именем
	greeting := fmt.Sprintf("👋 Добро пожаловать в бота управления подписками!\n\nПривет, %s!\n\nЭтот бот позволяет администраторам управлять подписками и пользователями.\n\nВыберите действие из меню ниже:", firstName)

	msg := tgbotapi.NewMessage(chatID, greeting)
	msg.ReplyMarkup = keyboards.MainAdminKeyboard()
	b.API.Send(msg)

	// Сбрасываем состояние пользователя
	if len(userID) > 0 {
		b.setUserState(userID[0], types.StateIdle)
	}
}

// showMainMenuEdit показывает главное меню через редактирование сообщения
func (b *Bot) showMainMenuEdit(chatID int64, messageID int, userID int64) {
	var firstName string

	// Получаем имя пользователя из Telegram API
	user, err := b.API.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: userID}})
	if err == nil {
		firstName = user.FirstName
	}

	// Если не удалось получить имя, используем по умолчанию
	if firstName == "" {
		firstName = "Администратор"
	}

	// Всегда показываем приветствие с именем
	greeting := fmt.Sprintf("👋 Добро пожаловать в бота управления подписками!\n\nПривет, %s!\n\nЭтот бот позволяет администраторам управлять подписками и пользователями.\n\nВыберите действие из меню ниже:", firstName)

	keyboard := keyboards.MainAdminKeyboard()
	b.editMessage(chatID, messageID, greeting, &keyboard)

	// Сбрасываем состояние пользователя
	b.setUserState(userID, types.StateIdle)
}

// showSubscriptionManagement показывает меню управления подписками
func (b *Bot) showSubscriptionManagement(chatID int64) {
	b.showMenu(chatID, 0, "subscriptions")
}

// showSubscriptionManagementEdit показывает меню управления подписками через редактирование
func (b *Bot) showSubscriptionManagementEdit(chatID int64, messageID int) {
	b.showMenu(chatID, messageID, "subscriptions")
}

// showUserManagement показывает меню управления пользователями
func (b *Bot) showUserManagement(chatID int64) {
	b.showMenu(chatID, 0, "users")
}

// showUserManagementEdit показывает меню управления пользователями через редактирование
func (b *Bot) showUserManagementEdit(chatID int64, messageID int) {
	b.showMenu(chatID, messageID, "users")
}

// showGlobalSettings показывает глобальные настройки
func (b *Bot) showGlobalSettings(chatID int64) {
	b.showMenu(chatID, 0, "global_settings")
}

// showGlobalSettingsEdit показывает глобальные настройки через редактирование
func (b *Bot) showGlobalSettingsEdit(chatID int64, messageID int) {
	b.showMenu(chatID, messageID, "global_settings")
}

// showAnalyticsEdit показывает меню аналитики через редактирование
func (b *Bot) showAnalyticsEdit(chatID int64, messageID int) {
	b.showMenu(chatID, messageID, "analytics")
}

// showMenu универсальная функция для показа меню (отправка или редактирование)
func (b *Bot) showMenu(chatID int64, messageID int, menuType string) {
	var text string
	var keyboard tgbotapi.InlineKeyboardMarkup

	switch menuType {
	case "subscriptions":
		text = "📝 Управление подписками\n\nВыберите действие:"
		keyboard = keyboards.SubscriptionManagementKeyboard()
	case "users":
		text = "👥 Управление пользователями\n\nВыберите действие:"
		keyboard = keyboards.UserManagementKeyboard()
	case "global_settings":
		// Получаем настройки для глобальных настроек
		settings, err := b.Context.APIClient.GetGlobalSettings()
		if err != nil {
			text = "⚙️ Глобальные настройки\n\nНастройки еще не созданы.\n\nСоздать глобальные настройки?"
			keyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📝 Создать настройки", "edit_global_markup"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu_edit"),
				),
			)
		} else {
			text = fmt.Sprintf("⚙️ Глобальные настройки\n\nТекущая глобальная надбавка: %.2f%%\n\nВыберите действие:", settings.GlobalMarkupPercent)
			keyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📝 Изменить надбавку", "edit_global_markup"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "global_settings_edit"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu_edit"),
				),
			)
		}
	case "analytics":
		text = "📊 Аналитика и статистика\n\nВыберите тип отчета:"
		keyboard = keyboards.AnalyticsKeyboard()
	default:
		text = "❌ Неизвестный тип меню"
		keyboard = keyboards.MainAdminKeyboard()
	}

	if messageID != 0 {
		// Редактируем существующее сообщение
		b.editMessage(chatID, messageID, text, &keyboard)
	} else {
		// Отправляем новое сообщение
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		sent, err := b.API.Send(msg)
		if err == nil {
			// Сохраняем ID сообщения для дальнейшего редактирования
			userState := b.getUserState(chatID)
			userState.CurrentMessageID = sent.MessageID
			userState.CurrentChatID = chatID
		}
	}
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

// answerCallbackQuery отвечает на callback запрос
func (b *Bot) answerCallbackQuery(queryID, text string) {
	callback := tgbotapi.NewCallback(queryID, text)
	b.API.Request(callback)
}

// getOrCreateAdminUser находит админа в БД или создает его если он админ
func (b *Bot) getOrCreateAdminUser(tgID int64, firstName, lastName, username string) (*api.User, error) {
	log.Printf("Getting admin user info for TGID: %d", tgID)

	// Проверяем, является ли пользователь админом в списке ADMINS
	if !b.isAdmin(tgID) {
		log.Printf("User %d not in admin list", tgID)
		return nil, fmt.Errorf("user with TGID %d not found", tgID)
	}

	// Сначала пытаемся найти пользователя
	user, err := b.Context.APIClient.FindUserByTGID(tgID)
	if err == nil {
		log.Printf("Found existing user: %s (ID: %s, IsAdmin: %t)", user.Fullname, user.ID, user.IsAdmin)

		// Если пользователь уже админ, просто возвращаем его
		if user.IsAdmin {
			log.Printf("User %s is already admin, no upgrade needed", user.Fullname)
			return user, nil
		}

		// Если пользователь не админ, но находится в списке админов, повышаем его
		log.Printf("User %d is in admin list but not admin in DB – updating role", tgID)
		isAdmin := true
		updateReq := api.UpdateUserRequest{IsAdmin: &isAdmin}
		updated, upErr := b.Context.APIClient.UpdateUser(user.ID, updateReq)
		if upErr != nil {
			log.Printf("Failed to update user role: %v", upErr)
			return user, nil // возвращаем как есть, если не удалось обновить
		} else {
			log.Printf("Successfully promoted user %s to admin", updated.Fullname)
			return updated, nil
		}
	}

	log.Printf("User not found in database: %v", err)

	// Формируем ФИО из имени и фамилии
	fullname := firstName
	if lastName != "" {
		fullname = fmt.Sprintf("%s %s", firstName, lastName)
	}
	if fullname == "" {
		fullname = fmt.Sprintf("Admin %d", tgID)
	}

	log.Printf("Creating new admin user: %s (TGID: %d)", fullname, tgID)

	// Создаем админа автоматически
	req := api.CreateUserRequest{
		TGID:     tgID,
		Username: username,
		Fullname: fullname,
		IsAdmin:  true,
	}

	user, err = b.Context.APIClient.CreateUser(req)
	if err != nil {
		log.Printf("Failed to create user: %v", err)

		// Если ошибка 409 (дубликат TGID), пытаемся найти пользователя еще раз
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "duplicate") {
			log.Printf("Duplicate TGID error, trying to find user again")
			user, findErr := b.Context.APIClient.FindUserByTGID(tgID)
			if findErr == nil {
				log.Printf("Found existing user after duplicate error: %s (ID: %s, IsAdmin: %t)", user.Fullname, user.ID, user.IsAdmin)
				// Если пользователь найден, но не админ, делаем его админом
				if !user.IsAdmin {
					log.Printf("User found but not admin, upgrading...")
					isAdmin := true
					updateReq := api.UpdateUserRequest{IsAdmin: &isAdmin}
					updated, upErr := b.Context.APIClient.UpdateUser(user.ID, updateReq)
					if upErr != nil {
						log.Printf("Failed to update user role: %v", upErr)
						return user, nil
					} else {
						log.Printf("Successfully promoted user %s to admin", updated.Fullname)
						return updated, nil
					}
				}
				return user, nil
			}
			log.Printf("Still couldn't find user after duplicate error: %v", findErr)
		}

		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("Successfully created admin user: %s (ID: %s)", user.Fullname, user.ID)
	return user, nil
}

// handleCurrencySelection обрабатывает выбор валюты
func (b *Bot) handleCurrencySelection(query *tgbotapi.CallbackQuery) {
	userState := b.getUserState(query.From.ID)
	currency := strings.TrimPrefix(query.Data, "currency_")

	log.Printf("User %d selected currency: %s in state: %s", query.From.ID, currency, userState.State)

	switch userState.State {
	case types.StateAwaitingSubscriptionCurrency:
		// Создание новой подписки
		if userState.SubscriptionData == nil {
			log.Printf("ERROR: SubscriptionData is nil for user %d", query.From.ID)
			text := "❌ Ошибка: данные подписки не найдены."
			keyboard := keyboards.CreateSuccessKeyboard("manage_subscriptions")
			b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
			return
		}

		log.Printf("Setting currency for subscription. Before: %+v", userState.SubscriptionData)
		userState.SubscriptionData.BaseCurrency = currency
		userState.State = types.StateAwaitingSubscriptionPeriod
		log.Printf("Setting currency for subscription. After: %+v", userState.SubscriptionData)

		text := fmt.Sprintf("📝 Создание новой подписки\n\n**Шаг 4/4:** Введите период списания в днях\n\n✅ Название: %s\n✅ Цена: %.2f %s\n\n*Например: 30 для ежемесячной подписки*",
			userState.SubscriptionData.ServiceName, userState.SubscriptionData.BasePrice, currency)
		keyboard := keyboards.CreateProcessKeyboard("process")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)

	case types.StateEditingSubscriptionCurrency:
		// Редактирование существующей подписки
		if userState.EditData == nil {
			text := "❌ Ошибка: данные для редактирования не найдены."
			keyboard := keyboards.CreateSuccessKeyboard("manage_subscriptions")
			b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
			return
		}

		// Обновляем подписку
		req := api.UpdateSubscriptionRequest{
			BaseCurrency: &currency,
		}

		subscription, err := b.Context.APIClient.UpdateSubscription(userState.EditData.EntityID, req)
		if err != nil {
			text := fmt.Sprintf("❌ Ошибка при обновлении подписки: %v", err)
			keyboard := keyboards.CreateSuccessKeyboard("manage_subscriptions")
			b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
			return
		}

		// Сбрасываем состояние
		userState.State = types.StateIdle
		userState.EditData = nil

		text := fmt.Sprintf("✅ Валюта подписки обновлена: %s", subscription.BaseCurrency)
		keyboard := keyboards.CreateSuccessKeyboard(fmt.Sprintf("edit_sub_%s", subscription.ID))
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)

	default:
		text := "❌ Неожиданное действие."
		keyboard := keyboards.CreateSuccessKeyboard("main_menu")
		b.editMessage(userState.CurrentChatID, userState.CurrentMessageID, text, &keyboard)
	}
}

// editMessage редактирует существующее сообщение
func (b *Bot) editMessage(chatID int64, messageID int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	if keyboard != nil {
		edit.ReplyMarkup = keyboard
	}
	b.API.Send(edit)
}

// getAdminUser получает информацию об админе из Telegram и создает/находит в системе
func (b *Bot) getAdminUser(adminUserID int64) (*api.User, error) {
	log.Printf("Getting admin user info for TGID: %d", adminUserID)

	// Получаем информацию о пользователе из Telegram
	user, err := b.API.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: adminUserID}})
	if err != nil {
		log.Printf("Failed to get user info from Telegram: %v", err)
		return nil, fmt.Errorf("ошибка при получении информации о пользователе: %w", err)
	}

	log.Printf("Got Telegram user info: %s %s (@%s)", user.FirstName, user.LastName, user.UserName)

	// Ищем или создаем админ пользователя в системе
	adminUser, err := b.getOrCreateAdminUser(adminUserID, user.FirstName, user.LastName, user.UserName)
	if err != nil {
		log.Printf("Failed to get/create admin user: %v", err)
		return nil, fmt.Errorf("ошибка при поиске/создании админ пользователя: %w", err)
	}

	log.Printf("Successfully got admin user: %s (ID: %s)", adminUser.Fullname, adminUser.ID)
	return adminUser, nil
}
