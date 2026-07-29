# FoodFlow Restaurant Studio

The Restaurant Studio is an independently built and deployed administration app. It manages restaurants, menu categories, menu items, and item availability without coupling its release lifecycle to the customer storefront.

## Local development

```sh
npm install
npm run dev
```

The app runs at `http://localhost:8081`. Vite proxies Sales API requests to port `3000` and Auth API requests to port `6000`. Copy `.env.example` to `.env.local` only when the services use different addresses.

## Production build

```sh
npm ci
npm run build
```

From the repository root, `make admin-frontend` builds the standalone container image.

## Staging access

The `staging-admin` Cloud Run service is intentionally IAM-protected while the backend's development token bootstrap remains enabled. Use an authenticated local proxy to open it:

```sh
gcloud run services proxy staging-admin \
  --project project-da1b7a1f-5ac7-474f-a6c \
  --region asia-southeast1 \
  --port 8081
```

Then open `http://localhost:8081`. Do not grant `allUsers` the Cloud Run Invoker role until production-grade admin authentication is implemented.
