# Implementation Plan - Phases 6-8 (Final)

This document contains the final phases: Frontend Integration, Testing & Deployment, and Monitoring.

---

## Phase 6: Frontend Integration

**Status:** ⬜ Not Started  
**Directory:** `api/frontends/food-flow-online-hub/src/`  
**Estimated Time:** 3-4 days

### 6.1 Install Stripe.js

#### Tasks
- [ ] Install Stripe packages: `npm install @stripe/stripe-js @stripe/react-stripe-js`
- [ ] Configure environment variables
- [ ] Create Stripe context provider

```bash
npm install @stripe/stripe-js @stripe/react-stripe-js
```

```env
# .env
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_...
VITE_API_URL=http://localhost:3000
```

---

### 6.2 Order Service

#### Tasks
- [ ] Create `src/services/orderService.ts`
- [ ] Implement API calls for orders
- [ ] Add TypeScript types
- [ ] Add error handling

```typescript
// File: src/services/orderService.ts

import { api } from './api';

export interface CreateOrderRequest {
  restaurantId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  orderType: 'delivery' | 'pickup';
  paymentMethod: 'creditCard' | 'payAtLocation';
  items: OrderItemRequest[];
  deliveryAddress?: DeliveryAddressRequest;
  specialInstructions?: string;
}

export interface OrderItemRequest {
  menuItemId: string;
  quantity: number;
  specialInstructions?: string;
}

export interface DeliveryAddressRequest {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  deliveryInstructions?: string;
}

export interface Order {
  id: string;
  restaurantId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  orderType: 'delivery' | 'pickup';
  orderStatus: string;
  paymentStatus: string;
  paymentMethod: string;
  subtotal: number;
  deliveryFee: number;
  tax: number;
  total: number;
  specialInstructions?: string;
  stripePaymentIntentId?: string;
  items: OrderItem[];
  deliveryAddress?: DeliveryAddress;
  dateCreated: string;
  dateUpdated: string;
}

export interface OrderItem {
  id: string;
  menuItemId: string;
  menuItemName: string;
  menuItemPrice: number;
  quantity: number;
  specialInstructions?: string;
}

export interface DeliveryAddress {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  deliveryInstructions?: string;
}

export interface PaymentIntentResponse {
  clientSecret: string;
  orderId: string;
  amount: number;
  currency: string;
}

export const orderService = {
  // Create a new order
  createOrder: async (request: CreateOrderRequest): Promise<Order> => {
    const response = await api.post('/v1/orders', request);
    return response.data;
  },

  // Get order by ID
  getOrder: async (orderId: string): Promise<Order> => {
    const response = await api.get(`/v1/orders/${orderId}`);
    return response.data;
  },

  // Create payment intent
  createPaymentIntent: async (orderId: string): Promise<PaymentIntentResponse> => {
    const response = await api.post(`/v1/orders/${orderId}/payment/intent`);
    return response.data;
  },

  // Confirm payment
  confirmPayment: async (orderId: string): Promise<void> => {
    await api.post(`/v1/orders/${orderId}/payment/confirm`);
  },
};
```

---

### 6.3 Stripe Payment Component

#### Tasks
- [ ] Create `src/components/StripePaymentForm.tsx`
- [ ] Integrate Stripe Elements
- [ ] Handle payment submission
- [ ] Add loading and error states

```typescript
// File: src/components/StripePaymentForm.tsx

import React, { useState } from 'react';
import {
  PaymentElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';

interface StripePaymentFormProps {
  orderId: string;
  onSuccess: () => void;
  onError: (error: Error) => void;
}

export const StripePaymentForm: React.FC<StripePaymentFormProps> = ({
  orderId,
  onSuccess,
  onError,
}) => {
  const stripe = useStripe();
  const elements = useElements();
  const [isProcessing, setIsProcessing] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!stripe || !elements) {
      return;
    }

    setIsProcessing(true);

    try {
      const { error } = await stripe.confirmPayment({
        elements,
        confirmParams: {
          return_url: `${window.location.origin}/order-confirmation/${orderId}`,
        },
        redirect: 'if_required',
      });

      if (error) {
        toast.error(error.message || 'Payment failed');
        onError(new Error(error.message));
      } else {
        toast.success('Payment successful!');
        onSuccess();
      }
    } catch (err) {
      const error = err as Error;
      toast.error('Payment failed. Please try again.');
      onError(error);
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <PaymentElement />
      <Button
        type="submit"
        disabled={!stripe || isProcessing}
        className="w-full"
      >
        {isProcessing ? 'Processing...' : 'Pay Now'}
      </Button>
    </form>
  );
};
```

---

### 6.4 Update Checkout Pages

#### Tasks
- [ ] Update `CheckoutDesktop.tsx` with real order creation
- [ ] Update `CheckoutMobile.tsx` with real order creation
- [ ] Integrate Stripe payment form
- [ ] Handle both credit card and pay-at-location
- [ ] Add error handling and loading states

```typescript
// File: src/pages/CheckoutDesktop.tsx (key changes)

import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
import { orderService } from '@/services/orderService';
import { StripePaymentForm } from '@/components/StripePaymentForm';

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

const CheckoutDesktop: React.FC = () => {
  const [clientSecret, setClientSecret] = useState<string | null>(null);
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null);
  
  // ... existing state ...

  const onSubmitCustomerInfo = async (data: any) => {
    setIsSubmitting(true);
    
    try {
      // Create order
      const order = await orderService.createOrder({
        restaurantId,
        customerName: data.name,
        customerEmail: data.email,
        customerPhone: data.phone,
        orderType,
        paymentMethod,
        items: items.map(item => ({
          menuItemId: item.menuItem.id,
          quantity: item.quantity,
          specialInstructions: item.specialInstructions,
        })),
        deliveryAddress: orderType === 'delivery' ? {
          street: data.street,
          city: data.city,
          state: data.state,
          postalCode: data.postalCode,
          deliveryInstructions: data.deliveryInstructions,
        } : undefined,
      });

      setCurrentOrder(order);

      if (paymentMethod === 'creditCard') {
        // Create payment intent
        const paymentIntent = await orderService.createPaymentIntent(order.id);
        setClientSecret(paymentIntent.clientSecret);
        setStep(2); // Move to payment step
      } else {
        // Pay at location - order complete
        clearCart();
        toast.success('Order placed successfully!');
        navigate(`/order-confirmation/${order.id}`);
      }
    } catch (error) {
      toast.error('Failed to create order. Please try again.');
      console.error(error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handlePaymentSuccess = async () => {
    if (!currentOrder) return;

    try {
      await orderService.confirmPayment(currentOrder.id);
      clearCart();
      navigate(`/order-confirmation/${currentOrder.id}`);
    } catch (error) {
      toast.error('Failed to confirm payment.');
      console.error(error);
    }
  };

  const handlePaymentError = (error: Error) => {
    console.error('Payment error:', error);
    // Keep user on payment page to retry
  };

  // Render payment step
  if (step === 2 && clientSecret && paymentMethod === 'creditCard') {
    return (
      <Elements stripe={stripePromise} options={{ clientSecret }}>
        <Layout>
          <div className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold mb-8">Complete Payment</h1>
            <StripePaymentForm
              orderId={currentOrder!.id}
              onSuccess={handlePaymentSuccess}
              onError={handlePaymentError}
            />
          </div>
        </Layout>
      </Elements>
    );
  }

  // ... rest of component ...
};
```

---

### 6.5 Order Confirmation Page

#### Tasks
- [ ] Create `src/pages/OrderConfirmation.tsx`
- [ ] Create `src/pages/OrderConfirmationMobile.tsx`
- [ ] Display order details
- [ ] Show order status
- [ ] Add routes to App.tsx

```typescript
// File: src/pages/OrderConfirmation.tsx

import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { orderService, Order } from '@/services/orderService';
import Layout from '@/components/Layout';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { CheckCircle } from 'lucide-react';

const OrderConfirmation: React.FC = () => {
  const { orderId } = useParams<{ orderId: string }>();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchOrder = async () => {
      if (!orderId) return;
      
      try {
        const data = await orderService.getOrder(orderId);
        setOrder(data);
      } catch (error) {
        console.error('Failed to fetch order:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchOrder();
  }, [orderId]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!order) {
    return <div>Order not found</div>;
  }

  return (
    <Layout>
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <div className="text-center mb-8">
          <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h1 className="text-3xl font-bold mb-2">Order Confirmed!</h1>
          <p className="text-gray-600">
            Order #{order.id.slice(0, 8)}
          </p>
        </div>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Order Details</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="font-semibold">Status:</span>
                <span className="capitalize">{order.orderStatus}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-semibold">Payment Status:</span>
                <span className="capitalize">{order.paymentStatus}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-semibold">Order Type:</span>
                <span className="capitalize">{order.orderType}</span>
              </div>
              
              <div className="border-t pt-4">
                <h3 className="font-semibold mb-2">Items:</h3>
                {order.items.map((item) => (
                  <div key={item.id} className="flex justify-between py-2">
                    <span>
                      {item.quantity}x {item.menuItemName}
                    </span>
                    <span>${(item.menuItemPrice * item.quantity).toFixed(2)}</span>
                  </div>
                ))}
              </div>

              <div className="border-t pt-4 space-y-2">
                <div className="flex justify-between">
                  <span>Subtotal:</span>
                  <span>${order.subtotal.toFixed(2)}</span>
                </div>
                {order.deliveryFee > 0 && (
                  <div className="flex justify-between">
                    <span>Delivery Fee:</span>
                    <span>${order.deliveryFee.toFixed(2)}</span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span>Tax:</span>
                  <span>${order.tax.toFixed(2)}</span>
                </div>
                <div className="flex justify-between font-bold text-lg border-t pt-2">
                  <span>Total:</span>
                  <span>${order.total.toFixed(2)}</span>
                </div>
              </div>

              {order.deliveryAddress && (
                <div className="border-t pt-4">
                  <h3 className="font-semibold mb-2">Delivery Address:</h3>
                  <p>{order.deliveryAddress.street}</p>
                  <p>
                    {order.deliveryAddress.city}, {order.deliveryAddress.state}{' '}
                    {order.deliveryAddress.postalCode}
                  </p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <div className="text-center text-gray-600">
          <p>A confirmation email has been sent to {order.customerEmail}</p>
        </div>
      </div>
    </Layout>
  );
};

export default OrderConfirmation;
```

---

**Phase 6 Summary:**
- ✅ Stripe.js integrated
- ✅ Order service implemented
- ✅ Payment form created
- ✅ Checkout pages updated
- ✅ Order confirmation page created
- ✅ Ready for testing

---

## Phase 7: Testing & Deployment

**Status:** ⬜ Not Started  
**Estimated Time:** 3-4 days

### 7.1 Backend Testing

#### Unit Tests
- [ ] Test orderbus business logic
- [ ] Test paymentbus with mock Stripe
- [ ] Test validation functions
- [ ] Test model conversions

```bash
# Run tests
go test ./business/domain/orderbus/...
go test ./business/domain/paymentbus/...
go test ./app/domain/orderapi/...
```

#### Integration Tests
- [ ] Test order creation API
- [ ] Test order query API
- [ ] Test payment intent creation
- [ ] Test webhook handling
- [ ] Test end-to-end order flow

---

### 7.2 Frontend Testing

#### Component Tests
- [ ] Test checkout pages
- [ ] Test payment form
- [ ] Test order service
- [ ] Test error handling

#### E2E Tests
- [ ] Test complete checkout flow
- [ ] Test payment with test cards
- [ ] Test order confirmation
- [ ] Test error scenarios

```bash
# Run tests
npm test
npm run test:e2e
```

---

### 7.3 Database Migration

#### Tasks
- [ ] Review migration scripts
- [ ] Test migration in local environment
- [ ] Test rollback script
- [ ] Run migration in staging
- [ ] Verify data integrity

```bash
# Run migration
make migrate-up

# Rollback if needed
make migrate-down
```

---

### 7.4 Environment Configuration

#### Backend Configuration
```yaml
# infra/k8s/base/stripe-secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: stripe-secrets
  namespace: sales-system
type: Opaque
stringData:
  secret-key: "sk_test_..."
  publishable-key: "pk_test_..."
  webhook-secret: "whsec_..."
```

#### Frontend Configuration
```env
# Production .env
VITE_STRIPE_PUBLISHABLE_KEY=pk_live_...
VITE_API_URL=https://api.foodflow.com
```

---

### 7.5 Deployment Steps

#### Tasks
- [ ] Build backend services
- [ ] Build frontend
- [ ] Run migrations in production
- [ ] Deploy to staging
- [ ] Run smoke tests
- [ ] Deploy to production
- [ ] Monitor deployment

```bash
# Build and deploy
make all
make kind-load
kubectl apply -f infra/k8s/base/
kubectl rollout status deployment/sales -n sales-system
kubectl rollout status deployment/frontend -n sales-system
```

---

**Phase 7 Summary:**
- ✅ Backend tests written and passing
- ✅ Frontend tests written and passing
- ✅ Migration tested and verified
- ✅ Environment configured
- ✅ Successfully deployed
- ✅ Ready for monitoring

---

## Phase 8: Monitoring & Error Handling

**Status:** ⬜ Not Started  
**Estimated Time:** 2-3 days

### 8.1 Logging

#### Tasks
- [ ] Add structured logging for orders
- [ ] Log payment events
- [ ] Log webhook events
- [ ] Log errors with context
- [ ] Set up log aggregation

#### Log Categories
```go
// Order lifecycle
log.Info(ctx, "order_created", "order_id", order.ID, "total", order.Total)
log.Info(ctx, "payment_intent_created", "order_id", order.ID, "intent_id", pi.ID)
log.Info(ctx, "payment_succeeded", "order_id", order.ID)

// Errors
log.Error(ctx, "order_creation_failed", "err", err)
log.Error(ctx, "payment_failed", "order_id", order.ID, "err", err)
log.Error(ctx, "webhook_verification_failed", "err", err)
```

---

### 8.2 Metrics & Monitoring

#### Tasks
- [ ] Track order creation rate
- [ ] Track payment success/failure rate
- [ ] Track average order value
- [ ] Track API response times
- [ ] Set up dashboards

#### Key Metrics
```
# Orders
- orders_created_total
- orders_by_status{status="pending|confirmed|preparing|ready|completed|cancelled"}
- order_value_distribution
- average_order_processing_time

# Payments
- payments_successful_total
- payments_failed_total
- payment_failures_by_reason{reason="declined|insufficient_funds|expired_card"}
- average_payment_processing_time

# Performance
- api_request_duration_seconds
- database_query_duration_seconds
- stripe_api_duration_seconds
```

---

### 8.3 Alerting

#### Tasks
- [ ] Alert on high payment failure rate (>5%)
- [ ] Alert on order creation failures
- [ ] Alert on webhook processing delays
- [ ] Alert on database connection issues
- [ ] Alert on Stripe API errors

#### Alert Rules
```yaml
# High payment failure rate
- alert: HighPaymentFailureRate
  expr: rate(payments_failed_total[5m]) > 0.05
  annotations:
    summary: "High payment failure rate detected"

# Order creation failures
- alert: OrderCreationFailures
  expr: rate(order_creation_errors_total[5m]) > 3
  annotations:
    summary: "Multiple order creation failures detected"

# Webhook processing delays
- alert: WebhookProcessingDelay
  expr: webhook_processing_duration_seconds > 5
  annotations:
    summary: "Webhook processing taking too long"
```

---

### 8.4 Error Handling Strategies

#### Common Scenarios

| Scenario | Handling | Recovery |
|----------|----------|----------|
| Payment declined | Show user-friendly error | Allow retry with different card |
| Network timeout | Check order status | Prevent duplicate orders |
| Menu item unavailable | Validate before payment | Show error, update cart |
| Webhook failure | Retry with exponential backoff | Alert after 3 failures |
| Duplicate order | Use idempotency keys | Return existing order |
| Stripe API down | Circuit breaker pattern | Queue orders, process later |

---

### 8.5 Health Checks

#### Tasks
- [ ] Implement health check endpoint
- [ ] Check database connectivity
- [ ] Check Stripe API status
- [ ] Monitor service dependencies

```go
// Health check endpoint
GET /health

Response:
{
  "status": "healthy",
  "database": "connected",
  "stripe": "available",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

---

**Phase 8 Summary:**
- ✅ Comprehensive logging implemented
- ✅ Metrics and monitoring configured
- ✅ Alerting rules set up
- ✅ Error handling strategies documented
- ✅ Health checks operational
- ✅ System production-ready

---

## Overall Progress Summary

### Completion Status

```
Phase 1: Business Domain        ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
Phase 2: Application Domain     ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
Phase 3: Database & Store       ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
Phase 4: REST API               ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
Phase 5: Stripe Integration     ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
Phase 6: Frontend               ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
Phase 7: Testing & Deployment   ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
Phase 8: Monitoring             ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%

Total Progress:                 ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
```

### Estimated Timeline

| Phase | Duration | Dependencies |
|-------|----------|--------------|
| Phase 1 | 2-3 days | None |
| Phase 2 | 2 days | Phase 1 |
| Phase 3 | 1-2 days | Phase 1 |
| Phase 4 | 3-4 days | Phase 1, 2, 3 |
| Phase 5 | 3-4 days | Phase 1, 4 |
| Phase 6 | 3-4 days | Phase 4, 5 |
| Phase 7 | 3-4 days | Phase 1-6 |
| Phase 8 | 2-3 days | Phase 7 |
| **Total** | **19-28 days** | |

---

## Next Steps

### Immediate Actions
1. ✅ Review complete implementation plan
2. ⬜ Set up Stripe test account
3. ⬜ Create feature branch: `git checkout -b feature/ordering-system`
4. ⬜ Begin Phase 1: Implement orderbus domain models
5. ⬜ Commit and test incrementally

### Development Approach
- Follow TDD where applicable
- Commit frequently with descriptive messages
- Update this plan as we progress
- Check off tasks as completed
- Test each phase before moving to next

---

**Plan Complete!**  
**Ready to begin implementation.**  
**Last Updated:** November 6, 2025
