# food-flow
Backend for the food flow, restaurant management POS 

## Kubernetes development

The local Kubernetes deployment is managed by the Helm chart in `infra/helm/food-flow`. Create the Kind cluster and deploy it with:

```sh
make dev-up
make build dev-load dev-apply
```

The chart defaults in `infra/helm/food-flow/values.yaml` use normal pod networking, password-authenticated Postgres, and no bundled observability services. The local-only Kind overrides are in `infra/helm/food-flow/values-kind.yaml`; the Kind cluster configuration remains in `infra/k8s/dev/kind-config.yaml`.

Set `POSTGRES_PASSWORD` before running `make dev-apply`; the Make target stores it in a Kubernetes Secret rather than in the chart values.

The local Kind deployment also applies `values-kind.yaml`, which enables the database migration and seed init container. Staging and production must use the default chart values (or environment-specific files that keep `sales.migration.enabled: false`).

Kind also starts the observability stack: Grafana at `http://localhost:3100`, Prometheus at `http://localhost:9090`, Tempo at `http://localhost:3200`, and Loki/Promtail for pod logs. Grafana provisions dashboards for application health, Go runtime behaviour, HTTP requests, latency, errors, logs, and Tempo trace links.



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
