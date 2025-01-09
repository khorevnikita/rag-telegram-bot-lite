package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/controllers"
	"html/template"
	"net/http"
)

type Router struct {
	bot                    *tb.Bot
	Engine                 *gin.Engine
	CloudPaymentController controllers.CloudPaymentController
}

func NewRouter(bot *tb.Bot) Router {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"toJSON": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	})

	router.LoadHTMLGlob("./templates/*")

	return Router{
		CloudPaymentController: controllers.NewCloudPaymentController(bot),
		Engine:                 router,
		bot:                    bot,
	}
}

func (r Router) Listen() {

	api := r.Engine.Group("/api")
	{
		// Health check на /api/
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "I`m working..."})
		})

		// Создаём группу для webhooks
		webhooks := api.Group("/webhooks")
		{
			// Создаём вложенную группу для cloud-payments
			cloudPayments := webhooks.Group("/cloud-payments")
			{
				// Определяем маршрут с параметром :method
				cloudPayments.POST("/:method", r.CloudPaymentController.WebhookHandler)
			}
		}

		api.GET("checkout/cloud-payments", r.CloudPaymentController.CheckoutPage)
	}

	_ = r.Engine.Run(":8080")
}
