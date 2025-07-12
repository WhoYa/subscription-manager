package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/WhoYa/subscription-manager/internal/bot/types"
)

// Общие функции для работы с сообщениями

// sendErrorMessage отправляет сообщение об ошибке с кнопкой "Назад"
func (b *Bot) sendErrorMessage(chatID int64, messageID int, err error, backAction string) {
	text := fmt.Sprintf(MessageError, err)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonBack, backAction),
		),
	)
	if messageID > 0 {
		b.editMessage(chatID, messageID, text, &keyboard)
	} else {
		b.sendMessageWithKeyboard(chatID, text, &keyboard)
	}
}

// sendMessageWithKeyboard отправляет сообщение с клавиатурой
func (b *Bot) sendMessageWithKeyboard(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	b.API.Send(msg)
}

// sendSimpleMessage отправляет простое текстовое сообщение
func (b *Bot) sendSimpleMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.API.Send(msg)
}

// sendMessage отправляет простое текстовое сообщение.
// Это алиас для backward compatibility. Предпочитайте sendSimpleMessage или sendMessageWithKeyboard.
func (b *Bot) sendMessage(chatID int64, text string) {
	b.sendSimpleMessage(chatID, text)
}

// Функции для работы с пользовательскими данными

// resetUserState сбрасывает состояние пользователя
func (b *Bot) resetUserState(userID int64) {
	userState := b.getUserState(userID)
	userState.State = types.StateIdle
	userState.SubscriptionData = nil
	userState.UserCreateData = nil
	userState.EditData = nil
	userState.CurrentEntityID = ""
}

// Функции валидации

// validateFloat64 проверяет и парсит число с плавающей точкой
func validateFloat64(input string, minValue float64) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("пустое значение")
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат числа")
	}

	if value < minValue {
		return 0, fmt.Errorf("значение должно быть больше или равно %.2f", minValue)
	}

	return value, nil
}

// validateInt проверяет и парсит целое число
func validateInt(input string, minValue int) (int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("пустое значение")
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("неверный формат числа")
	}

	if value < minValue {
		return 0, fmt.Errorf("значение должно быть больше или равно %d", minValue)
	}

	return value, nil
}

// validateInt64 проверяет и парсит int64
func validateInt64(input string, minValue int64) (int64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("пустое значение")
	}

	value, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат числа")
	}

	if value < minValue {
		return 0, fmt.Errorf("значение должно быть больше или равно %d", minValue)
	}

	return value, nil
}

// validateString проверяет строку на пустоту
func validateString(input string, canBeEmpty bool) (string, error) {
	input = strings.TrimSpace(input)
	if !canBeEmpty && input == "" {
		return "", fmt.Errorf("значение не может быть пустым")
	}
	return input, nil
}

// Функции для работы со статусами

// getSubscriptionStatus возвращает текст статуса подписки
func getSubscriptionStatus(isActive bool) string {
	if isActive {
		return StatusActive
	}
	return StatusInactive
}

// getUserRoleStatus возвращает текст роли пользователя
func getUserRoleStatus(isAdmin bool) string {
	if isAdmin {
		return StatusAdmin
	}
	return StatusRegular
}

// formatUsername форматирует username пользователя
func formatUsername(username string) string {
	if username == "" {
		return StatusNotSet
	}
	return "@" + username
}

// Функции для работы с callback данными

// Функции для работы с логированием

// logError логирует ошибку с дополнительным контекстом
func logError(context string, err error) {
	log.Printf("[ERROR] %s: %v", context, err)
}

// logInfo логирует информационное сообщение
func logInfo(context string, message string) {
	log.Printf("[INFO] %s: %s", context, message)
}

// Функции для работы с API ошибками

// handleAPIError обрабатывает ошибки API и возвращает пользовательское сообщение
func handleAPIError(err error, context string) string {
	logError(context, err)

	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "404") || strings.Contains(errStr, "not found"):
		return "Запрашиваемая информация не найдена"
	case strings.Contains(errStr, "409") || strings.Contains(errStr, "duplicate"):
		return "Данные уже существуют"
	case strings.Contains(errStr, "400") || strings.Contains(errStr, "bad request"):
		return "Неверные данные запроса"
	case strings.Contains(errStr, "403") || strings.Contains(errStr, "forbidden"):
		return "Недостаточно прав"
	case strings.Contains(errStr, "500") || strings.Contains(errStr, "internal server"):
		return "Внутренняя ошибка сервера"
	default:
		return "Неизвестная ошибка сервера"
	}
}
