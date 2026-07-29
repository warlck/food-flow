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
