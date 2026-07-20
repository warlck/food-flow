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

`values.yaml` is safe for staging and production: it disables the `migrate-seed` init container. `values-kind.yaml` enables it only for the local Kind workflow. The `stripe-secrets` Secret remains externally managed; use `make dev-stripe-secrets` to create it locally.

## Migrating an existing Kustomize deployment

For a local Kind environment, recreate the cluster (`make dev-down`, then `make dev-up`) before using the Helm workflow. This is the safest migration because the chart uses release-scoped resource names.

For a persistent cluster, first back up the database. Delete the old Deployments, Services, ConfigMaps, and StatefulSet, but retain the `database-data` PVC. Then install the chart with `--set database.persistence.existingClaim=database-data`; this reuses the existing database volume while avoiding Helm ownership conflicts. Ensure the database password Secret uses the existing database password.
