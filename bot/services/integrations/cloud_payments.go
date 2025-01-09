package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/models"
	"gorag-telegram-bot/utils"
	"io"
	"net/http"
	"net/url"
	"time"
)

type WebhookRequest struct {
	TransactionID     string  `json:"transaction_id"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	PaymentAmount     float64 `json:"payment_amount"`
	PaymentCurrency   string  `json:"payment_currency"`
	OperationType     string  `json:"operation_type"`
	AccountID         string  `json:"account_id"`
	SubscriptionID    string  `json:"subscription_id"`
	DateTime          string  `json:"date_time"`
	IpAddress         string  `json:"ip_address"`
	IpCountry         string  `json:"ip_country"`
	IpCity            string  `json:"ip_city"`
	IpRegion          string  `json:"ip_region"`
	CardID            string  `json:"card_id"`
	CardFirstSix      string  `json:"card_first_six"`
	CardLastFour      string  `json:"card_last_four"`
	CardType          string  `json:"card_type"`
	CardExpDate       string  `json:"card_exp_date"`
	Issuer            string  `json:"issuer"`
	IssuerBankCountry string  `json:"issuer_bank_country"`
	Description       string  `json:"description"`
	AuthCode          string  `json:"auth_code"`
	Token             string  `json:"token"`
	TestMode          bool    `json:"test_mode"`
	Status            string  `json:"status"`
	GatewayName       string  `json:"gateway_name"`
	TotalFee          float64 `json:"total_fee"`
}

type SubscriptionRequest struct {
	Token               string  `json:"token"`
	AccountId           string  `json:"accountId"`
	Description         string  `json:"description"`
	Email               string  `json:"email"`
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	RequireConfirmation bool    `json:"requireConfirmation"`
	StartDate           string  `json:"startDate"`
	Interval            string  `json:"interval"`
	Period              int     `json:"period"`
}

type SubscriptionModel struct {
	Id                     string  `json:"Id"`
	AccountId              string  `json:"AccountId"`
	Description            string  `json:"Description"`
	Email                  string  `json:"Email"`
	Amount                 float64 `json:"Amount"`
	Currency               string  `json:"Currency"`
	RequireConfirmation    bool    `json:"RequireConfirmation"`
	StartDateIso           string  `json:"StartDateIso"`
	Interval               string  `json:"Interval"`
	Period                 int     `json:"Period"`
	Status                 string  `json:"Status"`
	NextTransactionDateIso string  `json:"NextTransactionDateIso"`
}

type SubscriptionResponse struct {
	Model   SubscriptionModel `json:"Model"`
	Success bool              `json:"Success"`
	Message string            `json:"Message"`
}

type SubscriptionsListResponse struct {
	Model   []SubscriptionModel `json:"Model"`
	Success bool                `json:"Success"`
	Message string              `json:"Message"`
}

type CloudPaymentsService struct {
	endpoint string
}

func GetCloudPaymentsService() CloudPaymentsService {
	return CloudPaymentsService{
		endpoint: "https://api.cloudpayments.ru",
	}
}

func (s CloudPaymentsService) ParseWebhookRequest(payload string) (*WebhookRequest, error) {
	// Разбираем form-urlencoded
	values, err := url.ParseQuery(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// Преобразуем значения в структуру
	request := &WebhookRequest{
		TransactionID:     values.Get("TransactionId"),
		Amount:            utils.ParseFloat(values.Get("Amount")),
		Currency:          values.Get("Currency"),
		PaymentAmount:     utils.ParseFloat(values.Get("PaymentAmount")),
		PaymentCurrency:   values.Get("PaymentCurrency"),
		OperationType:     values.Get("OperationType"),
		AccountID:         values.Get("AccountId"),
		SubscriptionID:    values.Get("SubscriptionId"),
		DateTime:          values.Get("DateTime"),
		IpAddress:         values.Get("IpAddress"),
		IpCountry:         values.Get("IpCountry"),
		IpCity:            utils.DecodeURLEncoding(values.Get("IpCity")),
		IpRegion:          utils.DecodeURLEncoding(values.Get("IpRegion")),
		CardID:            values.Get("CardId"),
		CardFirstSix:      values.Get("CardFirstSix"),
		CardLastFour:      values.Get("CardLastFour"),
		CardType:          values.Get("CardType"),
		CardExpDate:       values.Get("CardExpDate"),
		Issuer:            values.Get("Issuer"),
		IssuerBankCountry: values.Get("IssuerBankCountry"),
		Description:       utils.DecodeURLEncoding(values.Get("Description")),
		AuthCode:          values.Get("AuthCode"),
		Token:             values.Get("Token"),
		TestMode:          values.Get("TestMode") == "1",
		Status:            values.Get("Status"),
		GatewayName:       values.Get("GatewayName"),
		TotalFee:          utils.ParseFloat(values.Get("TotalFee")),
	}

	return request, nil
}

func (s CloudPaymentsService) CreateSubscription(token string, amount float64, uid uuid.UUID) (*string, error) {
	existingSubscriptions, err := s.GetSubscriptions(uid)
	if err != nil {
		return nil, err
	}

	for _, subscription := range existingSubscriptions.Model {
		if subscription.Status == "Active" && subscription.Amount == amount {
			return &subscription.Id, nil
		}
	}

	now := time.Now()
	nextPaymentDate := now.AddDate(0, config.AppConfig.Modules.Billing.Providers.CloudPayments.Period, 0)
	requestBody := SubscriptionRequest{
		Token:       token,
		AccountId:   uid.String(),
		Description: "Ежемесячная подписка на сервис example.com",
		//Email:               "user@example.com",
		Amount:              amount,
		Currency:            "RUB",
		RequireConfirmation: false,
		//StartDate:           "2021-11-02T21:00:00",
		StartDate: nextPaymentDate.Format("2006-01-02T15:04:05"),
		Interval:  config.AppConfig.Modules.Billing.Providers.CloudPayments.PeriodUnit,
		Period:    config.AppConfig.Modules.Billing.Providers.CloudPayments.Period,
	}

	body, err := s.apiRequest("/subscriptions/create", requestBody)
	if err != nil {
		return nil, err
	}

	var response SubscriptionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	spew.Dump(response)

	return &response.Model.Id, nil
}

func (s CloudPaymentsService) GetSubscriptions(uid uuid.UUID) (*SubscriptionsListResponse, error) {
	body, err := s.apiRequest("/subscriptions/find", map[string]string{
		"accountId": uid.String(),
	})
	if err != nil {
		return nil, err
	}

	var response SubscriptionsListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	spew.Dump("CP Subs list", response)

	return &response, nil
}

func (s CloudPaymentsService) CancelSubscription(user *models.User) error {
	existingSubscriptions, err := s.GetSubscriptions(user.ID)
	if err != nil {
		return err
	}

	for _, subscription := range existingSubscriptions.Model {
		if subscription.Status == "Active" {
			body, err := s.apiRequest("/subscriptions/cancel", map[string]string{
				"Id": subscription.Id,
			})
			spew.Dump("cancel", body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s CloudPaymentsService) apiRequest(path string, requestBody any) ([]byte, error) {
	// Сериализация тела запроса в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request body: %w", err)
	}

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", s.endpoint+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(config.AppConfig.Modules.Billing.Providers.CloudPayments.PublicKey, config.AppConfig.Modules.Billing.Providers.CloudPayments.SecretKey) // Замените на ваши ключи CloudPayments

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Чтение и разбор ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
