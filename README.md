# food-flow
Backend for the food flow, restaurant management POS 



Set environment variables:



# Backend
export SALES_STRIPE_SECRET_KEY="sk_test_..."
export SALES_STRIPE_WEBHOOK_SECRET="whsec_..."

# Frontend (.env)
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_...
VITE_API_URL=http://localhost:3000


Configure Stripe webhook in Dashboard:
Endpoint: https://your-domain/v1/webhooks/stripe
Events: payment_intent.succeeded, payment_intent.payment_failed
Test with Stripe test cards:

Success: 4242 4242 4242 4242
Decline: 4000 0000 0000 0002