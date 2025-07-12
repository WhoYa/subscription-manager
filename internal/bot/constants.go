package bot

// Константы для текстов сообщений
const (
	// Общие сообщения
	MessageError            = "❌ Произошла ошибка: %v"
	MessageUnknownCommand   = "Неизвестная команда. Используйте /help для получения справки."
	MessageUnknownAction    = "❌ Неизвестное действие."
	MessageNoAccess         = "❌ У вас нет прав доступа к этому боту."
	MessageUseStart         = "Используйте /start для начала работы или /help для получения справки."
	MessageInvalidCallback  = "❌ Неверный формат callback."
	MessageDataNotFound     = "❌ Ошибка: данные не найдены."
	MessageUnexpectedAction = "❌ Неожиданное действие."

	// Сообщения создания подписки
	MessageSubscriptionCreateStart  = "📝 Создание новой подписки\n\n**Шаг 1/4:** Введите название сервиса\n\n*Например: Netflix, Spotify*"
	MessageSubscriptionNameEmpty    = "📝 Создание новой подписки\n\n**Шаг 1/4:** Введите название сервиса\n\n❌ Название сервиса не может быть пустым. Попробуйте снова:\n\n*Например: Netflix, Spotify*"
	MessageSubscriptionPriceStep    = "📝 Создание новой подписки\n\n**Шаг 2/4:** Введите базовую стоимость\n\n✅ Название: %s\n\n*Например: 9.99*"
	MessageSubscriptionPriceError   = "📝 Создание новой подписки\n\n**Шаг 2/4:** Введите базовую стоимость\n\n✅ Название: %s\n\n❌ Неверный формат цены. Введите число больше 0:\n\n*Например: 9.99*"
	MessageSubscriptionCurrencyStep = "📝 Создание новой подписки\n\n**Шаг 3/4:** Выберите валюту\n\n✅ Название: %s\n✅ Цена: %.2f"
	MessageSubscriptionPeriodStep   = "📝 Создание новой подписки\n\n**Шаг 4/4:** Введите период списания в днях\n\n✅ Название: %s\n✅ Цена: %.2f %s\n\n*Например: 30 для ежемесячной подписки*"
	MessageSubscriptionPeriodError  = "📝 Создание новой подписки\n\n**Шаг 4/4:** Введите период списания в днях\n\n✅ Название: %s\n✅ Цена: %.2f %s\n\n❌ Неверный формат периода. Введите целое число больше 0:\n\n*Например: 30 для ежемесячной подписки*"
	MessageSubscriptionConfirm      = "📝 Подтверждение создания подписки\n\n🏷️ Сервис: %s\n💰 Цена: %.2f %s\n📅 Период: %d дней\n\n❓ Создать подписку?"
	MessageSubscriptionCreated      = "✅ Подписка успешно создана!\n\n🏷️ Сервис: %s\n💰 Цена: %.2f %s\n📅 Период: %d дней\n🆔 ID: %s"
	MessageSubscriptionCreateError  = "❌ Ошибка при создании подписки: %v"

	// Сообщения создания пользователя
	MessageUserCreateStart   = "👤 Создание нового пользователя\n\n**Шаг 1/3:** Введите ФИО пользователя\n\n*Например: Иван Петров*"
	MessageUserFullnameEmpty = "👤 Создание нового пользователя\n\n**Шаг 1/3:** Введите ФИО пользователя\n\n❌ ФИО не может быть пустым. Попробуйте снова:\n\n*Например: Иван Петров*"
	MessageUserTGIDStep      = "👤 Создание нового пользователя\n\n**Шаг 2/3:** Введите Telegram ID пользователя\n\n✅ ФИО: %s\n\n*Например: 123456789*"
	MessageUserTGIDError     = "👤 Создание нового пользователя\n\n**Шаг 2/3:** Введите Telegram ID пользователя\n\n✅ ФИО: %s\n\n❌ Неверный формат Telegram ID. Введите положительное число:\n\n*Например: 123456789*"
	MessageUserUsernameStep  = "👤 Создание нового пользователя\n\n**Шаг 3/3:** Введите username пользователя\n\n✅ ФИО: %s\n✅ Telegram ID: %d\n\n*Без символа @, можно оставить пустым*"
	MessageUserConfirm       = "👤 Подтверждение создания пользователя\n\n👤 ФИО: %s\n🆔 Telegram ID: %d\n📝 Username: %s\n\n❓ Создать пользователя?"
	MessageUserCreated       = "✅ Пользователь успешно создан!\n\n👤 ФИО: %s\n🆔 Telegram ID: %d\n📝 Username: %s\n🆔 ID: %s"
	MessageUserCreateError   = "❌ Ошибка при создании пользователя: %v"

	// Сообщения глобальных настроек
	MessageGlobalMarkupStart       = "📝 Изменение глобальной надбавки\n\nТекущая надбавка: %s\n\nВведите новое значение надбавки в процентах:\n\n*Например: 15.5*"
	MessageGlobalMarkupError       = "📝 Изменение глобальной надбавки\n\n❌ Неверный формат надбавки. Введите число больше или равное 0:\n\n*Например: 15.5*"
	MessageGlobalMarkupSet         = "✅ Глобальная надбавка установлена: %.2f%%"
	MessageGlobalMarkupSetWithNote = "✅ Глобальная надбавка установлена: %.2f%% (не удалось подтвердить)"

	// Сообщения редактирования
	MessageEditSubscriptionTitle = "📝 Редактирование подписок\n\nВыберите подписку для редактирования:"
	MessageEditSubscriptionEmpty = "📝 Редактирование подписок\n\n📭 Подписки не найдены.\n\nСоздайте первую подписку!"
	MessageEditSubscriptionMenu  = "📝 Редактирование подписки\n\n🏷️ Сервис: %s\n💰 Цена: %.2f %s\n📅 Период: %d дней\n📊 Статус: %s\n\nЧто хотите изменить?"
	MessageEditUserTitle         = "📝 Редактирование пользователей\n\nВыберите пользователя для редактирования:"
	MessageEditUserEmpty         = "📝 Редактирование пользователей\n\n📭 Пользователи не найдены.\n\nСоздайте первого пользователя!"
	MessageEditUserMenu          = "📝 Редактирование пользователя\n\n👤 ФИО: %s\n🆔 Telegram ID: %d\n📝 Username: %s\n🔑 Роль: %s\n\nЧто хотите изменить?"

	// Сообщения редактирования полей
	MessageEditSubscriptionNamePrompt   = "📝 Введите новое название подписки:"
	MessageEditSubscriptionPricePrompt  = "💰 Введите новую цену подписки:"
	MessageEditSubscriptionPeriodPrompt = "📅 Введите новый период подписки в днях:"
	MessageEditUserFullnamePrompt       = "👤 Введите новое ФИО пользователя:"
	MessageEditUserUsernamePrompt       = "📱 Введите новый username пользователя (без @, можно оставить пустым):"

	// Сообщения валидации при редактировании
	MessageEditNameEmpty     = "❌ Название не может быть пустым. Попробуйте снова:"
	MessageEditPriceError    = "❌ Неверный формат цены. Введите число больше 0:"
	MessageEditPeriodError   = "❌ Неверный формат периода. Введите целое число больше 0:"
	MessageEditFullnameEmpty = "❌ ФИО не может быть пустым. Попробуйте снова:"

	// Сообщения об успешном обновлении
	MessageSubscriptionNameUpdated     = "✅ Название подписки обновлено: %s"
	MessageSubscriptionPriceUpdated    = "✅ Цена подписки обновлена: %.2f %s"
	MessageSubscriptionPeriodUpdated   = "✅ Период подписки обновлен: %d дней"
	MessageSubscriptionCurrencyUpdated = "✅ Валюта подписки обновлена: %s"
	MessageUserFullnameUpdated         = "✅ ФИО пользователя обновлено: %s"
	MessageUserUsernameUpdated         = "✅ Username пользователя обновлен: %s"

	// Сообщения переключения статуса
	MessageSubscriptionActivated   = "✅ Подписка %s активирована"
	MessageSubscriptionDeactivated = "✅ Подписка %s деактивирована"
	MessageUserMadeAdmin           = "✅ Пользователь %s теперь администратор"
	MessageUserMadeRegular         = "✅ Пользователь %s теперь обычный пользователь"

	// Статусы
	StatusActive   = "✅ Активна"
	StatusInactive = "❌ Неактивна"
	StatusAdmin    = "👑 Администратор"
	StatusRegular  = "👤 Обычный пользователь"
	StatusNotSet   = "не указан"
	StatusUnknown  = "неизвестна"

	// Кнопки
	ButtonBack               = "◀️ Назад"
	ButtonCancel             = "❌ Отменить"
	ButtonCreate             = "✅ Создать"
	ButtonConfirm            = "✅ Подтвердить"
	ButtonCreateSubscription = "➕ Создать подписку"
	ButtonCreateUser         = "➕ Создать пользователя"
	ButtonEditName           = "📝 Название"
	ButtonEditPrice          = "💰 Цена"
	ButtonEditCurrency       = "💱 Валюта"
	ButtonEditPeriod         = "📅 Период"
	ButtonEditFullname       = "📝 ФИО"
	ButtonEditUsername       = "📱 Username"
	ButtonToggleStatus       = "🔄 Статус"
	ButtonToggleRole         = "🔑 Роль"
)
