#!/bin/sh
set -eu

PROJECT_ID=${GCP_PROJECT_ID:-project-da1b7a1f-5ac7-474f-a6c}
REGION=${GCP_REGION:-asia-southeast1}
REPO=${GCP_ARTIFACT_REPOSITORY:-food-flow-staging}
NETWORK=${GCP_NETWORK:-staging-vpc}
SUBNET=${GCP_SUBNET:-staging-subnet}
POSTGRES_HOST=${GCP_POSTGRES_HOST:-10.0.1.2}
POSTGRES_ZONE=${GCP_POSTGRES_ZONE:-asia-southeast1-b}
POSTGRES_VM=${GCP_POSTGRES_VM:-staging-postgres}
WORKLOAD_SERVICE_ACCOUNT=${GCP_WORKLOAD_SERVICE_ACCOUNT:-staging-workload@${PROJECT_ID}.iam.gserviceaccount.com}
MIGRATION_JOB=${GCP_MIGRATION_JOB:-staging-db-migrate}
IMAGE_PLATFORM=${GCP_IMAGE_PLATFORM:-linux/amd64}
RUN_MIGRATION=${GCP_RUN_MIGRATION:-true}
REGISTRY="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}"
VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "latest")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
AUTH_IMAGE="${REGISTRY}/auth:${VERSION}"
SALES_IMAGE="${REGISTRY}/sales:${VERSION}"
STOREFRONT_IMAGE="${REGISTRY}/storefront:${VERSION}"

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
        --set-secrets "SALES_DB_PASSWORD=staging-postgres-password:latest" \
        --max-retries 0 \
        --task-timeout 5m \
        --execute-now \
        --wait \
        --quiet
fi

echo "=> Deploying Auth Service..."
gcloud run deploy staging-auth \
    --project "${PROJECT_ID}" \
    --image "${AUTH_IMAGE}" \
    --region "${REGION}" \
    --port 6000 \
    --ingress internal \
    --network "${NETWORK}" \
    --subnet "${SUBNET}" \
    --vpc-egress private-ranges-only \
    --service-account "${WORKLOAD_SERVICE_ACCOUNT}" \
    --update-env-vars "AUTH_DB_HOST=${POSTGRES_HOST},AUTH_DB_DISABLE_TLS=true" \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

echo "=> Deploying Sales Service..."
gcloud run deploy staging-sales \
    --project "${PROJECT_ID}" \
    --image "${SALES_IMAGE}" \
    --region "${REGION}" \
    --port 3000 \
    --ingress internal \
    --network "${NETWORK}" \
    --subnet "${SUBNET}" \
    --vpc-egress private-ranges-only \
    --service-account "${WORKLOAD_SERVICE_ACCOUNT}" \
    --update-env-vars "SALES_DB_HOST=${POSTGRES_HOST},SALES_DB_DISABLE_TLS=true" \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

echo "=> Fetching Sales Service URL..."
SALES_URL=$(gcloud run services describe staging-sales --project "${PROJECT_ID}" --region "${REGION}" --format 'value(status.url)')
echo "Sales URL: ${SALES_URL}"

echo "=> Fetching Auth Service URL..."
AUTH_URL=$(gcloud run services describe staging-auth --project "${PROJECT_ID}" --region "${REGION}" --format 'value(status.url)')
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

echo "=> Deploying Storefront Service..."
gcloud run deploy staging-storefront \
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

STOREFRONT_URL=$(gcloud run services describe staging-storefront --project "${PROJECT_ID}" --region "${REGION}" --format 'value(status.url)')

echo "=> Verifying Storefront and private API proxies..."
curl --fail --show-error --silent --retry 6 --retry-all-errors --retry-delay 5 "${STOREFRONT_URL}/health"
curl --fail --show-error --silent --retry 6 --retry-all-errors --retry-delay 5 "${STOREFRONT_URL}/api/auth/v1/readiness"
curl --fail --show-error --silent --retry 6 --retry-all-errors --retry-delay 5 "${STOREFRONT_URL}/api/sales/v1/readiness"

echo "=> Deployment Complete!"
echo "Storefront is live at: ${STOREFRONT_URL}"
