package handlers

import (
	"github.com/google/uuid"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/models"
	"gorag-telegram-bot/services"
	"gorag-telegram-bot/services/integrations"
)

type BillingHandlers struct {
	Config    config.BillingConfig
	CPService integrations.CloudPaymentsService
}

func GetBillingHandlers() BillingHandlers {
	return BillingHandlers{
		Config:    config.AppConfig.Modules.Billing,
		CPService: integrations.GetCloudPaymentsService(),
	}
}

func (h BillingHandlers) RegisterCommands() {
	if h.Config.Enabled {
		// add some handlers
		core.Bot.Handle(&tb.InlineButton{Unique: "billing_subscribe"}, h.onSubscribe)
		core.Bot.Handle(&tb.InlineButton{Unique: "billing_unsubscribe"}, h.onUnsubscribe)
		core.Bot.Handle(&tb.InlineButton{Unique: "billing_unsubscribe_confirm"}, h.onUnsubscribeConfirm)
		core.Bot.Handle(&tb.InlineButton{Unique: "billing_unsubscribe_cancel"}, h.onUnsubscribeCancel)

		core.Bot.Handle(tb.OnCheckout, h.onCheckout)
		core.Bot.Handle(tb.OnPayment, h.onPaid)
	}
}

func (h BillingHandlers) onSubscribe(c tb.Context) error {
	_ = c.Respond()
	return services.SendSubscribeMessage(c)
}

func (h BillingHandlers) onUnsubscribe(c tb.Context) error {
	_ = c.Respond()
	return c.Send(h.Config.UnsubscribeConfirmationMessage, &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	}, &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{
				tb.InlineButton{
					Unique: "billing_unsubscribe_cancel",
					Text:   h.Config.UnsubscribeCancel,
				},
			},
			{
				tb.InlineButton{
					Unique: "billing_unsubscribe_confirm",
					Text:   h.Config.UnsubscribeConfirm,
				},
			},
		},
	})
}

func (h BillingHandlers) onUnsubscribeConfirm(c tb.Context) error {
	user, _ := services.GetUserByTelegramID(c.Sender().ID)
	_ = c.Respond()
	err := services.Unsubscribe(user)
	if err != nil {
		return c.Send(err.Error())
	}

	if h.Config.Providers.CloudPayments.Enabled {
		err = h.CPService.CancelSubscription(user)
		if err != nil {
			return c.Send(err.Error())
		}
	}

	return c.Send(config.AppConfig.Modules.Billing.ByeMessage)
}

func (h BillingHandlers) onUnsubscribeCancel(c tb.Context) error {
	_ = c.Respond()
	if config.AppConfig.Menu.Enabled {
		return services.SendMenuMessage(c.Sender().ID)
	}
	return services.SendAIModeDisclaimer(c.Sender().ID)
}

func (h BillingHandlers) onCheckout(c tb.Context) error {
	user, _ := services.GetUserByTelegramID(c.Sender().ID)
	_ = c.Respond()
	preCheckout := c.PreCheckoutQuery()

	// Проверка данных заказа
	if preCheckout.Payload != user.ID.String() {
		return c.Accept("Invalid payload. Please try again.")
	}
	return c.Accept()
}

func (h BillingHandlers) onPaid(c tb.Context) error {
	successPayment := c.Payment()
	_ = c.Respond()

	if successPayment.Total < config.AppConfig.Modules.Billing.Providers.YooKassa.Price {
		return c.Send(h.Config.NotEnoughMoneyNotification)
	}

	uid, err := uuid.Parse(successPayment.Payload)
	if err != nil {
		return c.Send(err.Error())
	}
	user, err := services.FindUser(uid)
	if err != nil {
		return c.Send(err.Error())
	}

	provider, err := services.GetPaymentProvider(models.ProviderYooKassa)
	if err != nil {
		return c.Send(err.Error())
	}

	method, err := services.SavePaymentMethod(user.ID, provider.ID, "", "", "00/00")
	if err != nil {
		return c.Send(err.Error())
	}

	amount := float64(successPayment.Total / 100)

	subscription, err := services.GrantSubscription(user, method, amount)
	if err != nil {
		return c.Send(err.Error())
	}

	_, err = services.SavePayment(subscription.ID, provider.ID, amount)
	if err != nil {
		return c.Send(err.Error())
	}

	_ = c.Send(h.Config.SubscriptionGrantedNotification)

	if config.AppConfig.Modules.Form.Enabled {
		if user.FormCompletedAt == nil {
			return services.SendStartFormMessage(user.TelegramID)

		}
	}

	if config.AppConfig.Menu.Enabled {
		return services.SendMenuMessage(user.TelegramID)
	}

	return nil
}
