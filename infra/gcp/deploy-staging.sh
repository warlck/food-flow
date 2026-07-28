#!/bin/sh
set -eu

PROJECT_ID=${GCP_PROJECT_ID:-project-da1b7a1f-5ac7-474f-a6c}
REGION=${GCP_REGION:-asia-southeast1}
REPO=${GCP_ARTIFACT_REPOSITORY:-food-flow-staging}
REGISTRY="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}"
VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "latest")

echo "=> Setting gcloud project to ${PROJECT_ID}"
gcloud config set project "${PROJECT_ID}"

echo "=> Authenticating Docker with Artifact Registry"
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet

echo "=> Building Auth Service..."
docker build -f infra/docker/dockerfile.auth -t "${REGISTRY}/auth:${VERSION}" --build-arg "BUILD_REF=${VERSION}" .
docker push "${REGISTRY}/auth:${VERSION}"

echo "=> Building Sales Service..."
docker build -f infra/docker/dockerfile.sales -t "${REGISTRY}/sales:${VERSION}" --build-arg "BUILD_REF=${VERSION}" .
docker push "${REGISTRY}/sales:${VERSION}"

echo "=> Deploying Auth Service..."
gcloud run deploy staging-auth \
    --image "${REGISTRY}/auth:${VERSION}" \
    --region "${REGION}" \
    --port 6000 \
    --ingress internal \
    --network staging-vpc \
    --subnet staging-subnet \
    --vpc-egress private-ranges-only \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

echo "=> Deploying Sales Service..."
gcloud run deploy staging-sales \
    --image "${REGISTRY}/sales:${VERSION}" \
    --region "${REGION}" \
    --port 3000 \
    --ingress internal \
    --network staging-vpc \
    --subnet staging-subnet \
    --vpc-egress private-ranges-only \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

echo "=> Fetching Sales Service URL..."
SALES_URL=$(gcloud run services describe staging-sales --region "${REGION}" --format 'value(status.url)')
echo "Sales URL: ${SALES_URL}"

echo "=> Fetching Auth Service URL..."
AUTH_URL=$(gcloud run services describe staging-auth --region "${REGION}" --format 'value(status.url)')
echo "Auth URL: ${AUTH_URL}"

echo "=> Building Storefront Service..."
docker build -f infra/docker/dockerfile.frontend -t "${REGISTRY}/storefront:${VERSION}" \
    --build-arg "BUILD_REF=${VERSION}" \
    --build-arg VITE_API_URL="" \
    --build-arg VITE_STRIPE_PUBLISHABLE_KEY="dummy" \
    .
docker push "${REGISTRY}/storefront:${VERSION}"

echo "=> Deploying Storefront Service..."
gcloud run deploy staging-storefront \
    --image "${REGISTRY}/storefront:${VERSION}" \
    --region "${REGION}" \
    --port 8080 \
    --ingress all \
    --network staging-vpc \
    --subnet staging-subnet \
    --vpc-egress private-ranges-only \
    --set-env-vars "SALES_API_URL=${SALES_URL},AUTH_API_URL=${AUTH_URL}" \
    --min-instances 0 \
    --max-instances 1 \
    --cpu-throttling \
    --quiet

echo "=> Deployment Complete!"
echo "Storefront is live at: $(gcloud run services describe staging-storefront --region "${REGION}" --format 'value(status.url)')"
