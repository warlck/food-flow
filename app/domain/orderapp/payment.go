package orderapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/paymentintent"
	"github.com/stripe/stripe-go/v80/webhook"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/foundation/web"
)

// =============================================================================
// Payment Intent Response

// PaymentIntentResponse represents the response when creating a payment intent.
type PaymentIntentResponse struct {
	ClientSecret string  `json:"clientSecret"`
	OrderID      string  `json:"orderId"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
}

// Encode implements the encoder interface.
func (p PaymentIntentResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(p)
	return data, "application/json", err
}

// =============================================================================
// Payment Handlers

// createPaymentIntent creates a Stripe PaymentIntent for an order.
func (a *app) createPaymentIntent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	orderIDStr := web.Param(r, "order_id")

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return errs.NewFieldErrors("order_id", err)
	}

	// Get the order
	ord, err := a.orderBus.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("querybyid: orderID[%s]: %w", orderID, err)
	}

	// Check if payment already exists
	if ord.StripePaymentIntentID != "" {
		// Return existing payment intent
		pi, err := paymentintent.Get(ord.StripePaymentIntentID, nil)
		if err != nil {
			return errs.Newf(errs.Internal, "get payment intent: %s", err)
		}

		return web.Respond(ctx, w, PaymentIntentResponse{
			ClientSecret: pi.ClientSecret,
			OrderID:      ord.ID.String(),
			Amount:       ord.Total.Value(),
			Currency:     "usd",
		}, http.StatusOK)
	}

	// Only allow payment intent creation for credit card orders
	if ord.PaymentMethod != orderbus.PaymentMethodCreditCard {
		return errs.Newf(errs.InvalidArgument, "payment intent only allowed for credit card payments")
	}

	// Only allow payment intent creation for pending payment status
	if ord.PaymentStatus != orderbus.PaymentStatusPending {
		return errs.Newf(errs.InvalidArgument, "payment intent only allowed for orders with pending payment status")
	}

	// Convert total to cents (Stripe uses smallest currency unit)
	amountInCents := int64(ord.Total.Value() * 100)

	// Create Stripe PaymentIntent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String("usd"),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"order_id":      ord.ID.String(),
			"restaurant_id": ord.RestaurantID.String(),
			"customer_name": ord.CustomerName,
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return errs.Newf(errs.Internal, "create payment intent: %s", err)
	}

	// Update order with payment intent ID
	if err := a.orderBus.UpdateStripePaymentIntent(ctx, orderID, pi.ID); err != nil {
		return errs.Newf(errs.Internal, "update order with payment intent: %s", err)
	}

	return web.Respond(ctx, w, PaymentIntentResponse{
		ClientSecret: pi.ClientSecret,
		OrderID:      ord.ID.String(),
		Amount:       ord.Total.Value(),
		Currency:     "usd",
	}, http.StatusCreated)
}

// confirmPayment confirms that a payment has been completed.
func (a *app) confirmPayment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	orderIDStr := web.Param(r, "order_id")

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return errs.NewFieldErrors("order_id", err)
	}

	// Get the order
	ord, err := a.orderBus.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("querybyid: orderID[%s]: %w", orderID, err)
	}

	// Check if payment intent exists
	if ord.StripePaymentIntentID == "" {
		return errs.Newf(errs.InvalidArgument, "no payment intent found for order")
	}

	// Get the PaymentIntent from Stripe to verify status
	pi, err := paymentintent.Get(ord.StripePaymentIntentID, nil)
	if err != nil {
		return errs.Newf(errs.Internal, "get payment intent: %s", err)
	}

	// Update order status based on payment intent status
	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		if err := a.orderBus.UpdateStatus(ctx, orderID, orderbus.UpdateOrderStatus{
			PaymentStatus: orderbus.PaymentStatusPaid,
			OrderStatus:   orderbus.OrderStatusConfirmed,
		}); err != nil {
			return errs.Newf(errs.Internal, "update order status: %s", err)
		}

	case stripe.PaymentIntentStatusProcessing:
		// Payment is still processing, no action needed
		break

	case stripe.PaymentIntentStatusRequiresPaymentMethod,
		stripe.PaymentIntentStatusRequiresConfirmation,
		stripe.PaymentIntentStatusRequiresAction:
		// Payment is not complete
		return errs.Newf(errs.InvalidArgument, "payment not complete, status: %s", pi.Status)

	case stripe.PaymentIntentStatusCanceled:
		if err := a.orderBus.UpdateStatus(ctx, orderID, orderbus.UpdateOrderStatus{
			PaymentStatus: orderbus.PaymentStatusFailed,
		}); err != nil {
			return errs.Newf(errs.Internal, "update order status: %s", err)
		}
		return errs.Newf(errs.InvalidArgument, "payment was canceled")

	default:
		return errs.Newf(errs.Internal, "unknown payment intent status: %s", pi.Status)
	}

	// Return updated order
	ord, err = a.orderBus.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("querybyid: orderID[%s]: %w", orderID, err)
	}

	return web.Respond(ctx, w, ToAppOrder(ord), http.StatusOK)
}

// =============================================================================
// Webhook Handler

// webhookApp handles Stripe webhook events.
type webhookApp struct {
	orderBus      *orderbus.Business
	webhookSecret string
}

// newWebhookApp constructs a webhook handler.
func newWebhookApp(orderBus *orderbus.Business, webhookSecret string) *webhookApp {
	return &webhookApp{
		orderBus:      orderBus,
		webhookSecret: webhookSecret,
	}
}

// handleStripeWebhook processes Stripe webhook events.
func (a *webhookApp) handleStripeWebhook(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return errs.Newf(errs.InvalidArgument, "reading request body: %s", err)
	}

	// Verify webhook signature
	sigHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, a.webhookSecret)
	if err != nil {
		return errs.Newf(errs.InvalidArgument, "verifying webhook signature: %s", err)
	}

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return errs.Newf(errs.Internal, "parsing webhook payload: %s", err)
		}

		orderIDStr, ok := pi.Metadata["order_id"]
		if !ok {
			// Not our payment intent, ignore
			return web.Respond(ctx, w, map[string]string{"status": "ignored"}, http.StatusOK)
		}

		orderID, err := uuid.Parse(orderIDStr)
		if err != nil {
			return errs.Newf(errs.Internal, "parsing order_id from metadata: %s", err)
		}

		// Update order to paid
		if err := a.orderBus.UpdateStatus(ctx, orderID, orderbus.UpdateOrderStatus{
			PaymentStatus: orderbus.PaymentStatusPaid,
			OrderStatus:   orderbus.OrderStatusConfirmed,
		}); err != nil {
			return errs.Newf(errs.Internal, "updating order status: %s", err)
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return errs.Newf(errs.Internal, "parsing webhook payload: %s", err)
		}

		orderIDStr, ok := pi.Metadata["order_id"]
		if !ok {
			// Not our payment intent, ignore
			return web.Respond(ctx, w, map[string]string{"status": "ignored"}, http.StatusOK)
		}

		orderID, err := uuid.Parse(orderIDStr)
		if err != nil {
			return errs.Newf(errs.Internal, "parsing order_id from metadata: %s", err)
		}

		// Update order to failed
		if err := a.orderBus.UpdateStatus(ctx, orderID, orderbus.UpdateOrderStatus{
			PaymentStatus: orderbus.PaymentStatusFailed,
		}); err != nil {
			return errs.Newf(errs.Internal, "updating order status: %s", err)
		}

	default:
		// Ignore other event types
	}

	return web.Respond(ctx, w, map[string]string{"status": "received"}, http.StatusOK)
}
