package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client представляет HTTP клиент для взаимодействия с API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient создает новый API клиент
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Subscription структуры

type Subscription struct {
	ID           string  `json:"id"`
	ServiceName  string  `json:"service_name"`
	IconURL      string  `json:"icon_url,omitempty"`
	BasePrice    float64 `json:"base_price"`
	BaseCurrency string  `json:"base_currency"`
	IsActive     bool    `json:"is_active"`
	PeriodDays   int     `json:"period_days"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName  string  `json:"service_name"`
	BasePrice    float64 `json:"base_price"`
	BaseCurrency string  `json:"base_currency"`
	PeriodDays   int     `json:"period_days"`
}

// User структуры

type User struct {
	ID        string `json:"id"`
	TGID      int64  `json:"tgid"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	IsAdmin   bool   `json:"is_admin"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateUserRequest struct {
	TGID     int64  `json:"tgid"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	IsAdmin  bool   `json:"is_admin"`
}

// Global Settings структуры

type GlobalSettings struct {
	ID                  string  `json:"id"`
	GlobalMarkupPercent float64 `json:"global_markup_percent"`
	UpdatedAt           string  `json:"updated_at"`
	CreatedAt           string  `json:"created_at"`
}

type UpdateGlobalSettingsRequest struct {
	GlobalMarkupPercent float64 `json:"global_markup_percent"`
}

// CreateGlobalSettingsRequest структура для создания глобальных настроек
type CreateGlobalSettingsRequest struct {
	GlobalMarkupPercent float64 `json:"global_markup_percent"`
}

// Subscription методы

// CreateSubscription создает новую подписку
func (c *Client) CreateSubscription(req CreateSubscriptionRequest) (*Subscription, error) {
	url := fmt.Sprintf("%s/api/subscriptions", c.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("API: Creating subscription - URL: %s, Request: %s", url, string(body))

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("API: Failed to make request to %s: %v", url, err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("API: Request failed - Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var subscription Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		log.Printf("API: Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("API: Subscription created successfully: %+v", subscription)
	return &subscription, nil
}

// GetSubscriptions получает список подписок
func (c *Client) GetSubscriptions(limit, offset int) ([]Subscription, error) {
	url := fmt.Sprintf("%s/api/subscriptions?limit=%d&offset=%d", c.BaseURL, limit, offset)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var subscriptions []Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subscriptions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return subscriptions, nil
}

// GetSubscription получает подписку по ID
func (c *Client) GetSubscription(id string) (*Subscription, error) {
	url := fmt.Sprintf("%s/api/subscriptions/%s", c.BaseURL, id)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var subscription Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &subscription, nil
}

// User методы

// CreateUser создает нового пользователя
func (c *Client) CreateUser(req CreateUserRequest) (*User, error) {
	url := fmt.Sprintf("%s/api/users", c.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("API: Creating user - URL: %s, Request: %s", url, string(body))

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("API: Failed to make request to %s: %v", url, err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("API: Request failed - Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("API: Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("API: User created successfully: %+v", user)
	return &user, nil
}

// GetUsers получает список пользователей
func (c *Client) GetUsers(limit, offset int) ([]User, error) {
	url := fmt.Sprintf("%s/api/users?limit=%d&offset=%d", c.BaseURL, limit, offset)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return users, nil
}

// GetUser получает пользователя по ID
func (c *Client) GetUser(id string) (*User, error) {
	url := fmt.Sprintf("%s/api/users/%s", c.BaseURL, id)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// Global Settings методы

// GetGlobalSettings получает глобальные настройки
func (c *Client) GetGlobalSettings() (*GlobalSettings, error) {
	url := fmt.Sprintf("%s/api/settings", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var settings GlobalSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &settings, nil
}

// UpdateGlobalSettings обновляет глобальные настройки
func (c *Client) UpdateGlobalSettings(req UpdateGlobalSettingsRequest) (*GlobalSettings, error) {
	url := fmt.Sprintf("%s/api/settings", c.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var settings GlobalSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &settings, nil
}

// CreateGlobalSettings создает новые глобальные настройки
func (c *Client) CreateGlobalSettings(req CreateGlobalSettingsRequest) (*GlobalSettings, error) {
	url := fmt.Sprintf("%s/api/settings", c.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var settings GlobalSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &settings, nil
}

// Analytics методы (примеры основных)

// ProfitStats представляет статистику прибыли
type ProfitStats struct {
	TotalProfit   float64 `json:"total_profit"`
	PaymentCount  int     `json:"payment_count"`
	AverageProfit float64 `json:"average_profit"`
}

// GetTotalProfit получает общую статистику прибыли
func (c *Client) GetTotalProfit(adminUserID string) (*ProfitStats, error) {
	url := fmt.Sprintf("%s/api/admin/%s/profit/total", c.BaseURL, adminUserID)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var stats ProfitStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stats, nil
}

// GetMonthlyProfit получает статистику прибыли за месяц
func (c *Client) GetMonthlyProfit(adminUserID string, year, month int) (*ProfitStats, error) {
	url := fmt.Sprintf("%s/api/admin/%s/profit/monthly/%d/%d", c.BaseURL, adminUserID, year, month)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var stats ProfitStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stats, nil
}

// Profit Analytics структуры
type UserProfitStat struct {
	UserID       string  `json:"user_id"`
	UserFullname string  `json:"user_fullname"`
	TotalProfit  float64 `json:"total_profit"`
	PaymentCount int     `json:"payment_count"`
}

type SubscriptionProfitStat struct {
	SubscriptionID   string  `json:"subscription_id"`
	SubscriptionName string  `json:"subscription_name"`
	TotalProfit      float64 `json:"total_profit"`
	PaymentCount     int     `json:"payment_count"`
}

// Update requests
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Fullname *string `json:"fullname,omitempty"`
	IsAdmin  *bool   `json:"is_admin,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName  *string  `json:"service_name,omitempty"`
	IconURL      *string  `json:"icon_url,omitempty"`
	BasePrice    *float64 `json:"base_price,omitempty"`
	BaseCurrency *string  `json:"base_currency,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
	PeriodDays   *int     `json:"period_days,omitempty"`
}

// Utility функции

// IsAdminUser проверяет, является ли пользователь администратором
func (c *Client) IsAdminUser(userID string) (bool, error) {
	user, err := c.GetUser(userID)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

// FindUserByTGID находит пользователя по Telegram ID
func (c *Client) FindUserByTGID(tgid int64) (*User, error) {
	// Используем новый эндпоинт для поиска пользователя по TGID
	url := fmt.Sprintf("%s/api/users/tgid/%d", c.BaseURL, tgid)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user with TGID %d not found", tgid)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// GetUserProfitStats получает статистику прибыли по пользователям
func (c *Client) GetUserProfitStats(adminUserID, from, to string) ([]UserProfitStat, error) {
	url := fmt.Sprintf("%s/api/admin/%s/profit/users?from=%s&to=%s", c.BaseURL, adminUserID, from, to)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var stats []UserProfitStat
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return stats, nil
}

// GetSubscriptionProfitStats получает статистику прибыли по подпискам
func (c *Client) GetSubscriptionProfitStats(adminUserID, from, to string) ([]SubscriptionProfitStat, error) {
	url := fmt.Sprintf("%s/api/admin/%s/profit/subscriptions?from=%s&to=%s", c.BaseURL, adminUserID, from, to)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var stats []SubscriptionProfitStat
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return stats, nil
}

// UpdateUser обновляет пользователя
func (c *Client) UpdateUser(id string, req UpdateUserRequest) (*User, error) {
	url := fmt.Sprintf("%s/api/users/%s", c.BaseURL, id)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// UpdateSubscription обновляет подписку
func (c *Client) UpdateSubscription(id string, req UpdateSubscriptionRequest) (*Subscription, error) {
	url := fmt.Sprintf("%s/api/subscriptions/%s", c.BaseURL, id)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var subscription Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &subscription, nil
}
