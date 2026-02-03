package orderapp

import (
	"net/http"

	"github.com/stripe/stripe-go/v80"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build               string
	Log                 *logger.Logger
	AuthClient          *authclient.Client
	OrderBus            *orderbus.Business
	StripeSecretKey     string
	StripeWebhookSecret string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	// Initialize Stripe with the secret key
	if cfg.StripeSecretKey != "" {
		stripe.Key = cfg.StripeSecretKey
	}

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.OrderBus)

	// Public order creation (customers can create orders)
	// NOTE: This route is unauthenticated as the ordering needs to stay public
	app.HandleFunc(http.MethodPost, version, "/orders", api.create)

	// Order queries
	// NOTE: Order confirmation page fetches order by id without auth in dev.
	app.HandleFunc(http.MethodGet, version, "/orders", api.query, authen)
	app.HandleFunc(http.MethodGet, version, "/orders/{order_id}", api.queryByID)

	// Status updates (admin only)
	app.HandleFunc(http.MethodPatch, version, "/orders/{order_id}/status", api.updateStatus, authen, ruleAdmin)

	// Cancel order (admin only)
	app.HandleFunc(http.MethodPost, version, "/orders/{order_id}/cancel", api.cancel, authen, ruleAdmin)

	// Payment endpoints
	// NOTE: Checkout flow is anonymous, so these routes are unauthenticated.
	app.HandleFunc(http.MethodPost, version, "/orders/{order_id}/payment/intent", api.createPaymentIntent)
	app.HandleFunc(http.MethodPost, version, "/orders/{order_id}/payment/confirm", api.confirmPayment)

	// Stripe webhook (no authentication - Stripe signs the payload)
	if cfg.StripeWebhookSecret != "" {
		webhookAPI := newWebhookApp(cfg.OrderBus, cfg.StripeWebhookSecret)
		app.HandleFunc(http.MethodPost, version, "/webhooks/stripe", webhookAPI.handleStripeWebhook)
	}
}
