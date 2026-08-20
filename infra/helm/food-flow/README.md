# Food Flow Helm chart

The chart replaces the Kustomize bases and development overlays while keeping the local Kind configuration in `infra/k8s/dev/kind-config.yaml`.

Install the development stack after creating the Kind cluster and loading the application images. The database password is supplied through an externally managed Secret:

```sh
kubectl create namespace sales-system
kubectl create secret generic food-flow-database-credentials \
  --namespace sales-system \
  --from-literal=password="$POSTGRES_PASSWORD"
helm upgrade --install food-flow infra/helm/food-flow \
  --namespace sales-system \
  --values infra/helm/food-flow/values-kind.yaml \
  --wait
```

`values.yaml` is safe for staging and production: it uses normal pod networking, password-authenticated Postgres, and disables the bundled observability services and `migrate-seed` init container. `values-kind.yaml` enables the local-only network exposure, disposable database trust authentication, observability stack, and migration/seed workflow. The `stripe-secrets` Secret remains externally managed; use `make dev-stripe-secrets` to create it locally.

The auth service signing key is also externally managed. The chart wires `AUTH_AUTH_KEYS_ENV_VAR` from the `food-flow-auth-keys` Secret (key `keys_json`, value `{"key":"<kid>","pem":"..."}`) and sets `AUTH_AUTH_ACTIVE_KID` from `auth.activeKID` in `values.yaml` (`local-dev` for kind). The secret reference is optional: if the secret is absent, the auth service falls back to reading PEMs from `infra/keys/` (offline `go run` development only). Use `make dev-auth-keys` to generate a throwaway local key and install the secret into kind; `make dev-apply` runs it automatically. For staging and production the key lives in GCP Secret Manager and is injected by the deploy scripts — see `infra/gcp/README.md`.

The Kind profile also starts Prometheus, Grafana, Tempo, Loki, and Promtail. Grafana provisions the **Food Flow Overview**, **Food Flow Go Runtime**, and **Food Flow Requests, Errors & Logs** dashboards automatically. Open Grafana at `http://localhost:3100`; log entries expose a TraceID link that opens the corresponding trace in Tempo.

## Migrating an existing Kustomize deployment

For a local Kind environment, recreate the cluster (`make dev-down`, then `make dev-up`) before using the Helm workflow. This is the safest migration because the chart uses release-scoped resource names.

For a persistent cluster, first back up the database. Delete the old Deployments, Services, ConfigMaps, and StatefulSet, but retain the `database-data` PVC. Then install the chart with `--set database.persistence.existingClaim=database-data`; this reuses the existing database volume while avoiding Helm ownership conflicts. Ensure the database password Secret uses the existing database password.
