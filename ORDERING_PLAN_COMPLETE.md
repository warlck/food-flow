# Food Flow - Complete Ordering & Checkout Implementation Plan

**Version:** 1.0  
**Last Updated:** November 6, 2025  
**Status:** Ready for Implementation  
**Payment Provider:** Stripe

---

## 📋 Plan Overview

This document provides a complete roadmap for implementing an ordering/checkout system with Stripe payment processing for the Food Flow restaurant POS application.

### Implementation Order
1. Business Domain Models (orderbus)
2. Application Domain Layer (orderapp)
3. Database Schema & Migrations
4. REST API Implementation (orderapi)
5. Stripe Payment Integration (paymentbus)
6. Frontend Integration (React + Stripe.js)
7. Testing & Deployment
8. Monitoring & Error Handling

---

## 📁 Plan Documents

### Main Documents

1. **ORDERING_IMPLEMENTATION_PLAN.md** *(311 lines)*
   - Architecture overview
   - Payment flow diagrams
   - **Phase 1:** Business Domain Models (orderbus)
   
2. **ORDERING_PLAN_PHASE2-8.md** *(~1000+ lines)*
   - **Phase 2:** Application Domain Layer (orderapp)
   - **Phase 3:** Database Schema & Migrations
   - **Phase 4:** REST API Implementation (orderapi)
   - **Phase 5:** Stripe Payment Integration (paymentbus)

3. **ORDERING_PLAN_PHASE6-8.md** *(this file)*
   - **Phase 6:** Frontend Integration (React + Stripe.js)
   - **Phase 7:** Testing & Deployment
   - **Phase 8:** Monitoring & Error Handling

---

## 🏗️ Architecture Layers

```
┌─────────────────────────────────────────┐
│     Frontend (React + TypeScript)       │
│   - Checkout Pages                      │
│   - Stripe Payment Form                 │
│   - Order Confirmation                  │
└──────────────┬──────────────────────────┘
               │ HTTP/REST
┌──────────────▼──────────────────────────┐
│        REST API (orderapi)              │
│   - POST /v1/orders                     │
│   - GET /v1/orders/:id                  │
│   - PATCH /v1/orders/:id/status         │
│   - POST /v1/orders/:id/payment/intent  │
│   - POST /webhooks/stripe               │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    Application Layer (orderapp)         │
│   - Request/Response Models             │
│   - Validation                          │
│   - Model Conversions                   │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│     Business Layer (orderbus)           │
│   - Order Domain Models                 │
│   - Business Logic & Rules              │
│   - Core Interface                      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    Database Store (orderdb)             │
│   - PostgreSQL                          │
│   - 4 Tables: orders, order_items,      │
│     delivery_addresses, payment_txns    │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│    Payment Service (paymentbus)         │
│   - Stripe SDK Integration              │
│   - PaymentIntent Creation              │
│   - Webhook Handling                    │
└─────────────────────────────────────────┘
```

---

## 🔄 Payment Flow

```
Customer → Add Items to Cart
         ↓
         Checkout Page
         ↓
         Enter Customer Info & Delivery Address
         ↓
         Select Payment Method (Credit Card / Pay at Location)
         ↓
   ┌─────┴──────┐
   │            │
Credit Card   Pay at Location
   │            │
   │            ├──→ Backend: Create Order (status: pending)
   │            └──→ Order Confirmation Page
   │
   ├──→ Backend: Create Order (status: pending)
   ├──→ Stripe: Create PaymentIntent
   ├──→ Frontend: Display Stripe Payment Form
   ├──→ Customer: Enter Card Details
   ├──→ Stripe: Process Payment
   ├──→ Stripe Webhook: payment_intent.succeeded
   ├──→ Backend: Update Order (payment_status: paid)
   └──→ Order Confirmation Page
```

---

## 📊 Database Schema

### Tables
1. **orders** - Main order records
2. **order_items** - Line items for each order
3. **delivery_addresses** - Delivery information
4. **payment_transactions** - Payment history

### Key Relationships
```
orders (1) ←→ (N) order_items
orders (1) ←→ (1) delivery_addresses
orders (1) ←→ (N) payment_transactions
```

---

## 🔑 Key Features

### Business Features
- ✅ Order creation with multiple items
- ✅ Delivery and pickup support
- ✅ Credit card and pay-at-location payment methods
- ✅ Order status tracking (pending → confirmed → preparing → ready → completed)
- ✅ Payment status tracking (pending → processing → paid)
- ✅ Order cancellation
- ✅ Delivery address management
- ✅ Special instructions for orders and items

### Payment Features
- ✅ Stripe PaymentIntent API
- ✅ Secure payment processing
- ✅ Webhook event handling
- ✅ Payment confirmation
- ✅ Refund support
- ✅ Payment failure handling

### API Features
- ✅ RESTful endpoints
- ✅ JWT authentication
- ✅ Role-based authorization
- ✅ Input validation
- ✅ Error handling
- ✅ Pagination and filtering

---

## 📝 Implementation Checklist

### Phase 1: Business Domain Models ⬜
- [ ] Define order status constants
- [ ] Define payment status constants
- [ ] Create Order domain model
- [ ] Create OrderItem domain model
- [ ] Create DeliveryAddress domain model
- [ ] Define Core business interface
- [ ] Define Storer interface
- [ ] Create query filters
- [ ] Add validation logic

### Phase 2: Application Domain Layer ⬜
- [ ] Create request models
- [ ] Create response models
- [ ] Add validation tags
- [ ] Implement model conversions
- [ ] Add filter parsing

### Phase 3: Database Schema ⬜
- [ ] Write migration scripts (up/down)
- [ ] Create orders table
- [ ] Create order_items table
- [ ] Create delivery_addresses table
- [ ] Create payment_transactions table
- [ ] Add indexes for performance
- [ ] Implement Store (orderdb)
- [ ] Write database models
- [ ] Create conversion functions
- [ ] Test migrations locally

### Phase 4: REST API ⬜
- [ ] Define routes
- [ ] Implement createOrder handler
- [ ] Implement queryOrders handler
- [ ] Implement queryOrderByID handler
- [ ] Implement updateOrderStatus handler
- [ ] Implement cancelOrder handler
- [ ] Register routes in sales service
- [ ] Add authentication middleware
- [ ] Add authorization rules
- [ ] Write API tests

### Phase 5: Stripe Integration ⬜
- [ ] Create Stripe account
- [ ] Obtain API keys
- [ ] Install Stripe SDK
- [ ] Configure Kubernetes secrets
- [ ] Define Payment business interface
- [ ] Implement Stripe client
- [ ] Implement CreatePaymentIntent
- [ ] Implement ConfirmPayment
- [ ] Implement RefundPayment
- [ ] Implement HandleWebhook
- [ ] Create payment API handlers
- [ ] Configure Stripe webhook endpoint
- [ ] Test with Stripe test cards

### Phase 6: Frontend Integration ⬜
- [ ] Install Stripe.js and React Stripe Elements
- [ ] Configure environment variables
- [ ] Create order service (API client)
- [ ] Create TypeScript types
- [ ] Create StripePaymentForm component
- [ ] Update CheckoutDesktop with real API calls
- [ ] Update CheckoutMobile with real API calls
- [ ] Handle credit card payment flow
- [ ] Handle pay-at-location flow
- [ ] Create OrderConfirmation pages
- [ ] Add error handling
- [ ] Add loading states
- [ ] Test payment flow end-to-end

### Phase 7: Testing & Deployment ⬜
- [ ] Write orderbus unit tests
- [ ] Write paymentbus unit tests
- [ ] Write API integration tests
- [ ] Write frontend component tests
- [ ] Write E2E tests
- [ ] Run migrations in staging
- [ ] Deploy backend to staging
- [ ] Deploy frontend to staging
- [ ] Run smoke tests
- [ ] Deploy to production
- [ ] Monitor deployment

### Phase 8: Monitoring & Error Handling ⬜
- [ ] Implement structured logging
- [ ] Set up metrics collection
- [ ] Create monitoring dashboards
- [ ] Configure alerting rules
- [ ] Document error scenarios
- [ ] Implement health check endpoint
- [ ] Test failure scenarios
- [ ] Set up log aggregation

---

## 🛠️ Technology Stack

### Backend
- **Language:** Go 1.24.3
- **Framework:** Custom (based on existing architecture)
- **Database:** PostgreSQL
- **Payment:** Stripe SDK v80
- **Authentication:** JWT
- **Container:** Docker
- **Orchestration:** Kubernetes

### Frontend
- **Framework:** React 18
- **Language:** TypeScript
- **Build Tool:** Vite
- **Payment:** @stripe/stripe-js, @stripe/react-stripe-js
- **HTTP Client:** Axios
- **Styling:** Tailwind CSS

### Infrastructure
- **Container Registry:** Docker Hub / Private Registry
- **Deployment:** Kubernetes
- **CI/CD:** GitHub Actions (to be configured)
- **Monitoring:** Prometheus + Grafana
- **Logging:** Structured JSON logs

---

## 📁 File Structure

### Backend
```
business/domain/
  orderbus/
    orderbus.go          # Core business logic
    model.go             # Domain models
    filter.go            # Query filters
    order.go             # Ordering
    testutil.go          # Test helpers
    stores/
      orderdb/
        orderdb.go       # Database store
        model.go         # DB models
        
  paymentbus/
    paymentbus.go        # Payment service
    model.go             # Payment models
    stripe/
      stripe.go          # Stripe client

app/domain/
  orderapi/
    orderapi.go          # API handlers
    route.go             # Route definitions
    model.go             # Request/response models
    
business/domain/orderbus/stores/orderdb/
  migrate/sql/
    001_orders.sql       # Migration scripts
```

### Frontend
```
src/
  services/
    orderService.ts      # API client
    
  components/
    StripePaymentForm.tsx  # Payment form
    
  pages/
    CheckoutDesktop.tsx    # Desktop checkout
    CheckoutMobile.tsx     # Mobile checkout
    OrderConfirmation.tsx  # Confirmation page
```

---

## ⏱️ Timeline Estimate

| Phase | Duration | Start After |
|-------|----------|-------------|
| Phase 1: Business Domain | 2-3 days | - |
| Phase 2: Application Layer | 2 days | Phase 1 |
| Phase 3: Database | 1-2 days | Phase 1 |
| Phase 4: REST API | 3-4 days | Phase 1, 2, 3 |
| Phase 5: Stripe Integration | 3-4 days | Phase 1, 4 |
| Phase 6: Frontend | 3-4 days | Phase 4, 5 |
| Phase 7: Testing & Deployment | 3-4 days | Phase 1-6 |
| Phase 8: Monitoring | 2-3 days | Phase 7 |
| **Total** | **19-28 days** | |

### Critical Path
```
Phase 1 → Phase 2 → Phase 4 → Phase 5 → Phase 6 → Phase 7 → Phase 8
         ↓
       Phase 3 ↗
```

---

## 🚀 Getting Started

### Prerequisites
- [ ] Go 1.24.3 installed
- [ ] PostgreSQL 15+ running
- [ ] Node.js 18+ installed
- [ ] Stripe test account created
- [ ] Docker and Kubernetes access

### Initial Setup

1. **Create feature branch**
   ```bash
   git checkout -b feature/ordering-system
   ```

2. **Set up Stripe**
   - Sign up at https://stripe.com
   - Get test API keys from Dashboard
   - Note webhook secret

3. **Configure environment**
   ```bash
   # Backend
   export STRIPE_SECRET_KEY="sk_test_..."
   export STRIPE_WEBHOOK_SECRET="whsec_..."
   
   # Frontend
   echo "VITE_STRIPE_PUBLISHABLE_KEY=pk_test_..." > .env.local
   ```

4. **Install dependencies**
   ```bash
   # Backend
   go get github.com/stripe/stripe-go/v80
   
   # Frontend
   npm install @stripe/stripe-js @stripe/react-stripe-js
   ```

5. **Begin Phase 1**
   - Open `ORDERING_IMPLEMENTATION_PLAN.md`
   - Follow Phase 1 implementation steps
   - Commit frequently
   - Check off completed tasks

---

## 📚 Documentation References

### Internal
- [Main Plan - Phase 1](./ORDERING_IMPLEMENTATION_PLAN.md)
- [Phases 2-5](./ORDERING_PLAN_PHASE2-8.md)
- [Phases 6-8](./ORDERING_PLAN_PHASE6-8.md)

### External
- [Stripe API Docs](https://stripe.com/docs/api)
- [Stripe Payment Intents](https://stripe.com/docs/payments/payment-intents)
- [Stripe Webhooks](https://stripe.com/docs/webhooks)
- [Stripe Testing Cards](https://stripe.com/docs/testing)

---

## 🎯 Success Criteria

### Phase 1-5 Complete
- ✅ Orders can be created via API
- ✅ Orders are persisted in database
- ✅ Stripe PaymentIntents are created
- ✅ Payments can be confirmed
- ✅ Webhooks update order status

### Phase 6 Complete
- ✅ Users can complete checkout in frontend
- ✅ Credit card payments work end-to-end
- ✅ Pay-at-location orders work
- ✅ Order confirmation page displays correctly

### Phase 7 Complete
- ✅ All tests passing
- ✅ Successfully deployed to production
- ✅ No critical bugs in first week

### Phase 8 Complete
- ✅ Monitoring dashboards operational
- ✅ Alerts configured and tested
- ✅ Error scenarios documented
- ✅ Team trained on monitoring tools

---

## 🔒 Security Considerations

- ✅ Stripe API keys stored in Kubernetes secrets
- ✅ Webhook signature verification
- ✅ JWT authentication on API endpoints
- ✅ RBAC authorization
- ✅ Input validation on all requests
- ✅ SQL injection prevention (parameterized queries)
- ✅ HTTPS for all API communication
- ✅ PCI compliance (Stripe handles card data)

---

## 📞 Support & Questions

### During Implementation
- Review phase-specific documentation
- Check Stripe documentation for payment issues
- Review existing codebase patterns (userbus, checkapi)
- Ask team members for clarification

### Common Issues
- **Stripe test mode:** Ensure using test keys (pk_test_, sk_test_)
- **Webhook not firing:** Use Stripe CLI for local testing
- **Payment failing:** Use Stripe test cards
- **Database migration:** Test rollback before production

---

## ✅ Final Checklist Before Launch

- [ ] All 8 phases completed
- [ ] Tests passing (unit, integration, E2E)
- [ ] Code reviewed by team
- [ ] Database migrations tested
- [ ] Stripe production keys configured
- [ ] Webhook endpoint configured in Stripe Dashboard
- [ ] Monitoring dashboards set up
- [ ] Alerts configured
- [ ] Documentation updated
- [ ] Team trained on new features
- [ ] Rollback plan documented
- [ ] Production deployment successful
- [ ] Smoke tests passed
- [ ] First real order completed successfully 🎉

---

**Ready to begin implementation!**

**Start with:** `ORDERING_IMPLEMENTATION_PLAN.md` → Phase 1: Business Domain Models

**Good luck! 🚀**
