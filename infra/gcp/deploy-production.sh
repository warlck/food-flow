#!/bin/sh
set -eu

# Production deploy. Mirrors deploy-staging.sh with production hardening:
#   - its own project, Artifact Registry repo, secret, and workload SA
#     (nothing is shared with staging)
#   - the auth service runs min-instances 1 without CPU throttling because
#     the sales API delegates authentication to it on every call
#   - the admin frontend sits behind IAP (Google account gate) via the
#     direct Cloud Run IAP integration; the storefront stays public
#   - the auth service CORS allowlist is pinned to the admin origin
#
# Prerequisites:
#   - GCP_PROJECT_ID set to the production project
#   - GCP_AUTH_ACTIVE_KID set to the KID created by
#     infra/gcp/bootstrap-auth-secret.sh (run against the production project)
#   - the production Postgres VM running

PROJECT_ID=${GCP_PROJECT_ID:?GCP_PROJECT_ID is required (production project)}
REGION=${GCP_REGION:-asia-southeast1}
REPO=${GCP_ARTIFACT_REPOSITORY:-food-flow-production}
NETWORK=${GCP_NETWORK:-production-vpc}
SUBNET=${GCP_SUBNET:-production-subnet}
POSTGRES_HOST=${GCP_POSTGRES_HOST:?GCP_POSTGRES_HOST is required}
POSTGRES_ZONE=${GCP_POSTGRES_ZONE:-asia-southeast1-b}
POSTGRES_VM=${GCP_POSTGRES_VM:-production-postgres}
WORKLOAD_SERVICE_ACCOUNT=${GCP_WORKLOAD_SERVICE_ACCOUNT:-production-workload@${PROJECT_ID}.iam.gserviceaccount.com}
MIGRATION_JOB=${GCP_MIGRATION_JOB:-production-db-migrate}
IMAGE_PLATFORM=${GCP_IMAGE_PLATFORM:-linux/amd64}
RUN_MIGRATION=${GCP_RUN_MIGRATION:-true}
AUTH_SECRET_NAME=${GCP_AUTH_SECRET:-food-flow-auth-keys}
AUTH_ACTIVE_KID=${GCP_AUTH_ACTIVE_KID:-}
AUTH_SERVICE=${GCP_AUTH_SERVICE:-production-auth}
SALES_SERVICE=${GCP_SALES_SERVICE:-production-sales}
STOREFRONT_SERVICE=${GCP_STOREFRONT_SERVICE:-production-storefront}
ADMIN_SERVICE=${GCP_ADMIN_SERVICE:-production-admin}
REGISTRY="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}"
VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "latest")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
AUTH_IMAGE="${REGISTRY}/auth:${VERSION}"
SALES_IMAGE="${REGISTRY}/sales:${VERSION}"
STOREFRONT_IMAGE="${REGISTRY}/storefront:${VERSION}"
ADMIN_FRONTEND_IMAGE="${REGISTRY}/admin-frontend:${VERSION}"

for command_name in docker gcloud git curl; do
    if ! command -v "${command_name}" >/dev/null 2>&1; then
        echo "Missing required command: ${command_name}" >&2
        exit 1
    fi
done

if ! docker info >/dev/null 2>&1; then
    echo "Docker is not running. Start Docker Desktop and retry." >&2
    exit 1
fi

if [ -z "${AUTH_ACTIVE_KID}" ]; then
    echo "GCP_AUTH_ACTIVE_KID is not set. Run infra/gcp/bootstrap-auth-secret.sh create against the production project and retry." >&2
    exit 1
fi

vm_status=$(gcloud compute instances describe "${POSTGRES_VM}" \
    --project "${PROJECT_ID}" \
    --zone "${POSTGRES_ZONE}" \
    --format='value(status)')
if [ "${vm_status}" != "RUNNING" ]; then
    echo "${POSTGRES_VM} is ${vm_status:-unavailable}; start it before deploying." >&2
    exit 1
fi

echo "=> Setting gcloud project to ${PROJECT_ID}"
gcloud config set project "${PROJECT_ID}"

echo "=> Authenticating Docker with Artifact Registry"
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet

echo "=> Building Auth Service..."
docker buildx build \
    --platform "${IMAGE_PLATFORM}" \
    -f infra/docker/dockerfile.auth \
    -t "${AUTH_IMAGE}" \
    --build-arg "BUILD_REF=${VERSION}" \
    --build-arg "BUILD_DATE=${BUILD_DATE}" \
    --push \
    .

echo "=> Building Sales Service..."
docker buildx build \
    --platform "${IMAGE_PLATFORM}" \
    -f infra/docker/dockerfile.sales \
    -t "${SALES_IMAGE}" \
    --build-arg "BUILD_REF=${VERSION}" \
    --build-arg "BUILD_DATE=${BUILD_DATE}" \
    --push \
    .

if [ "${RUN_MIGRATION}" = "true" ]; then
    echo "=> Running database migrations..."
    gcloud run jobs deploy "${MIGRATION_JOB}" \
        --project "${PROJECT_ID}" \
        --image "${SALES_IMAGE}" \
        --region "${REGION}" \
        --command ./admin \
        --args migrate \
        --service-account "${WORKLOAD_SERVICE_ACCOUNT}" \
        --network "${NETWORK}" \
        --subnet "${SUBNET}" \
        --vpc-egress private-ranges-only \
        --set-env-vars "SALES_DB_HOST=${POSTGRES_HOST},SALES_DB_DISABLE_TLS=true" \
        --set-secrets "SALES_DB_PASSWORD=food-flow-db-password:latest" \
        --max-retries 0 \
        --task-timeout 5m \
        --execute-now \
        --wait \
        --quiet
fi

echo "=> Deploying Auth Service..."
gcloud run deploy "${AUTH_SERVICE}" \
    --project "${PROJECT_ID}" \
    --image "${AUTH_IMAGE}" \
    --region "${REGION}" \
    --port 6000 \
    --ingress internal \
    --network "${NETWORK}" \
    --subnet "${SUBNET}" \
    --vpc-egress private-ranges-only \
    --service-account "${WORKLOAD_SERVICE_ACCOUNT}" \
    --remove-env-vars AUTH_DB_DISABLE_TLS \
    --update-env-vars "AUTH_DB_HOST=${POSTGRES_HOST},AUTH_AUTH_ACTIVE_KID=${AUTH_ACTIVE_KID}" \
    --set-secrets "AUTH_AUTH_KEYS_ENV_VAR=${AUTH_SECRET_NAME}:latest,AUTH_DB_PASSWORD=food-flow-db-password:latest" \
    --min-instances 1 \
    --max-instances 1 \
    --quiet

echo "=> Deploying Sales Service..."
gcloud run deploy "${SALES_SERVICE}" \
    --project "${PROJECT_ID}" \
    --image "${SALES_IMAGE}" \
    --region "${REGION}" \
    --port 3000 \
    --ingress internal \
    --network "${NETWORK}" \
    --subnet "${SUBNET}" \
    --vpc-egress private-ranges-only \
    --service-account "${WORKLOAD_SERVICE_ACCOUNT}" \
    --update-env-vars "SALES_DB_HOST=${POSTGRES_HOST}" \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

echo "=> Fetching Sales Service URL..."
SALES_URL=$(gcloud run services describe "${SALES_SERVICE}" --project "${PROJECT_ID}" --region "${REGION}" --format 'value(status.url)')
echo "Sales URL: ${SALES_URL}"

echo "=> Fetching Auth Service URL..."
AUTH_URL=$(gcloud run services describe "${AUTH_SERVICE}" --project "${PROJECT_ID}" --region "${REGION}" --format 'value(status.url)')
echo "Auth URL: ${AUTH_URL}"

echo "=> Building Storefront Service..."
docker buildx build \
    --platform "${IMAGE_PLATFORM}" \
    -f infra/docker/dockerfile.frontend \
    -t "${STOREFRONT_IMAGE}" \
    --build-arg "BUILD_REF=${VERSION}" \
    --build-arg "BUILD_DATE=${BUILD_DATE}" \
    --build-arg VITE_API_URL="" \
    --build-arg "VITE_STRIPE_PUBLISHABLE_KEY=${VITE_STRIPE_PUBLISHABLE_KEY:-dummy}" \
    --push \
    .

echo "=> Deploying Storefront Service (public, no IAP)..."
gcloud run deploy "${STOREFRONT_SERVICE}" \
    --project "${PROJECT_ID}" \
    --image "${STOREFRONT_IMAGE}" \
    --region "${REGION}" \
    --port 8080 \
    --ingress all \
    --network "${NETWORK}" \
    --subnet "${SUBNET}" \
    --vpc-egress private-ranges-only \
    --service-account "${WORKLOAD_SERVICE_ACCOUNT}" \
    --update-env-vars "SALES_API_URL=${SALES_URL},AUTH_API_URL=${AUTH_URL}" \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

echo "=> Building Admin Frontend Service..."
docker buildx build \
    --platform "${IMAGE_PLATFORM}" \
    -f infra/docker/dockerfile.admin-frontend \
    -t "${ADMIN_FRONTEND_IMAGE}" \
    --build-arg "BUILD_REF=${VERSION}" \
    --build-arg "BUILD_DATE=${BUILD_DATE}" \
    --build-arg VITE_SALES_API_URL="" \
    --build-arg VITE_AUTH_API_URL="" \
    --push \
    .

echo "=> Deploying Admin Frontend Service..."
gcloud run deploy "${ADMIN_SERVICE}" \
    --project "${PROJECT_ID}" \
    --image "${ADMIN_FRONTEND_IMAGE}" \
    --region "${REGION}" \
    --port 8080 \
    --ingress all \
    --no-allow-unauthenticated \
    --network "${NETWORK}" \
    --subnet "${SUBNET}" \
    --vpc-egress private-ranges-only \
    --service-account "${WORKLOAD_SERVICE_ACCOUNT}" \
    --update-env-vars "SALES_API_URL=${SALES_URL},AUTH_API_URL=${AUTH_URL}" \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

STOREFRONT_URL=$(gcloud run services describe "${STOREFRONT_SERVICE}" --project "${PROJECT_ID}" --region "${REGION}" --format 'value(status.url)')
ADMIN_FRONTEND_URL=$(gcloud run services describe "${ADMIN_SERVICE}" --project "${PROJECT_ID}" --region "${REGION}" --format 'value(status.url)')

echo "=> Pinning auth CORS allowlist to the admin origin..."
gcloud run services update "${AUTH_SERVICE}" \
    --project "${PROJECT_ID}" \
    --region "${REGION}" \
    --update-env-vars "AUTH_WEB_CORS_ALLOWED_ORIGINS=${ADMIN_FRONTEND_URL}" \
    --quiet

echo "=> Verifying Storefront and private API proxies..."
curl --fail --show-error --silent --retry 6 --retry-all-errors --retry-delay 5 "${STOREFRONT_URL}/health"
curl --fail --show-error --silent --retry 6 --retry-all-errors --retry-delay 5 "${STOREFRONT_URL}/api/auth/v1/readiness"
curl --fail --show-error --silent --retry 6 --retry-all-errors --retry-delay 5 "${STOREFRONT_URL}/api/sales/v1/readiness"

# Verify the admin frontend before IAP is enabled: the identity token works
# against the Cloud Run invoker check. After IAP is on, programmatic checks
# need an ID token minted for the IAP OAuth client audience instead.
echo "=> Verifying Admin Frontend (pre-IAP, invoker identity token)..."
ADMIN_IDENTITY_TOKEN=$(gcloud auth print-identity-token)
curl --fail --show-error --silent --retry 6 --retry-all-errors --retry-delay 5 \
    --header "Authorization: Bearer ${ADMIN_IDENTITY_TOKEN}" \
    "${ADMIN_FRONTEND_URL}/health"

echo "=> Enabling IAP on the Admin Frontend..."
gcloud services enable iap.googleapis.com --project "${PROJECT_ID}" --quiet

PROJECT_NUMBER=$(gcloud projects describe "${PROJECT_ID}" --format='value(projectNumber)')
gcloud run services add-iam-policy-binding "${ADMIN_SERVICE}" \
    --project "${PROJECT_ID}" \
    --region "${REGION}" \
    --member="serviceAccount:service-${PROJECT_NUMBER}@gcp-sa-iap.iam.gserviceaccount.com" \
    --role="roles/run.invoker" \
    --quiet

gcloud beta run services update "${ADMIN_SERVICE}" \
    --project "${PROJECT_ID}" \
    --region "${REGION}" \
    --iap \
    --quiet

cat <<EOF

=> Deployment Complete!
Storefront is live at: ${STOREFRONT_URL}
Restaurant Studio is live at: ${ADMIN_FRONTEND_URL} (IAP-gated)

One-time IAP setup (console, once per project):
  1. Configure the OAuth consent screen (APIs & Services > OAuth consent screen).
  2. Grant each operator the IAP-secured Web App User role
     (roles/iap.httpsResourceAccessor) on the ${ADMIN_SERVICE} service.
  3. Programmatic checks against the gated admin frontend now need an ID
     token minted for the IAP OAuth client audience:
       gcloud auth print-identity-token --audiences=<IAP_CLIENT_ID>
EOF
