package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/models"
	"gorag-telegram-bot/services"
	"gorag-telegram-bot/services/integrations"
	"io"
	"net/http"
	"os"
)

type CloudPaymentController struct {
	billingConfig config.BillingConfig
	bot           *tb.Bot
	service       integrations.CloudPaymentsService
}

func NewCloudPaymentController(bot *tb.Bot) CloudPaymentController {
	return CloudPaymentController{
		billingConfig: config.AppConfig.Modules.Billing,
		bot:           bot,
		service:       integrations.GetCloudPaymentsService(),
	}
}

func (c *CloudPaymentController) CheckoutPage(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Query("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID пользователя введен некорректно",
		})
		return
	}

	user, err := services.FindUser(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	data := gin.H{
		"PublicID":    c.billingConfig.Providers.CloudPayments.PublicKey, // Замените на ваш Public ID
		"Description": fmt.Sprintf("Подписка на ежемесячный доступ к @%s", config.AppConfig.BotUsername),
		"Amount":      c.billingConfig.Providers.CloudPayments.Price, // Сумма подписки
		"Currency":    "RUB",
		"AccountID":   user.ID.String(),
		"Period":      c.billingConfig.Providers.CloudPayments.Period,
		"PeriodUnit":  c.billingConfig.Providers.CloudPayments.PeriodUnit,
		"Items": []gin.H{
			{
				"label":    "Ежемесячный платёж",
				"price":    c.billingConfig.Providers.CloudPayments.Price,
				"quantity": 1,
				"amount":   c.billingConfig.Providers.CloudPayments.Price,
				"vat":      0,
				"method":   0,
				"object":   0,
			},
		},
	}

	ctx.HTML(200, "cloud_payments_checkout.html", data)
}

func (c *CloudPaymentController) WebhookHandler(ctx *gin.Context) {
	method := ctx.Param("method")
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if method != "pay" && method != "confirm" {
		ctx.JSON(200, gin.H{
			"code": 0,
		})
		return
	}

	// Парсинг запроса
	request, err := c.service.ParseWebhookRequest(string(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing payload: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    0,
			"message": "Invalid form data",
		})
		return
	}

	userID, err := uuid.Parse(request.AccountID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Wrong Account ID: %v\n", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	}

	user, err := services.FindUser(userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Wrong Account ID: %v\n", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	}

	if request.Amount < config.AppConfig.Modules.Billing.Providers.CloudPayments.Price {
		_ = services.SendNotification(user, c.billingConfig.NotEnoughMoneyNotification)
		fmt.Fprintf(os.Stderr, "Small amount: %v\n", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	}

	err = c.handleNewPayment(user, request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	}

	c.notifyOnSuccess(user)

	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func (c *CloudPaymentController) handleNewPayment(user *models.User, request *integrations.WebhookRequest) error {
	pp, err := services.GetPaymentProvider(models.ProviderCloudPayments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannt find payment provider: %v\n", err)
		return err
	}

	pm, err := services.SavePaymentMethod(user.ID, pp.ID, request.Token, request.CardLastFour, request.CardExpDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannt save payment method: %v\n", err)
		return err
	}

	subscription, err := services.GrantSubscription(user, pm, request.Amount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error giving subscription: %v\n", err)

		return err
	}

	_, err = services.SavePayment(subscription.ID, pp.ID, request.Amount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving payment: %v\n", err)
		return err
	}

	_, err = c.service.CreateSubscription(pm.Token, subscription.Amount, user.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating subscription in CP: %v\n", err)
		return err
	}

	_ = services.SendNotification(user, c.billingConfig.SubscriptionGrantedNotification)
	return nil
}

func (c *CloudPaymentController) notifyOnSuccess(user *models.User) {
	_ = services.SendNotification(user, c.billingConfig.SubscriptionGrantedNotification)

	if config.AppConfig.Modules.Form.Enabled {
		if user.FormCompletedAt == nil {
			_ = services.SendStartFormMessage(user.TelegramID)
			return
		}
	}

	if config.AppConfig.Menu.Enabled {
		_ = services.SendMenuMessage(user.TelegramID)
	}
}
