#!/bin/sh
set -eu

# Provision or rotate the auth service JWT signing key in GCP Secret Manager.
#
# The stored payload matches what foundation/keystore LoadByJSON expects:
#
#     {"key": "<kid>", "pem": "<private PEM>"}
#
# Cloud Run injects it into the auth service as the AUTH_AUTH_KEYS_ENV_VAR
# environment variable via --set-secrets (see deploy-staging.sh).
#
# Usage:
#   GCP_PROJECT_ID=prj GCP_WORKLOAD_SERVICE_ACCOUNT=sa@prj.iam.gserviceaccount.com \
#       ./infra/gcp/bootstrap-auth-secret.sh create
#   GCP_PROJECT_ID=prj ./infra/gcp/bootstrap-auth-secret.sh rotate
#
# create  Provision the secret, add the first key version, and grant the
#         environment's workload service account secretAccessor.
# rotate  Add a new key version with a fresh KID. The secret and its IAM
#         binding must already exist. Follow the rotation runbook: keep the
#         previous version enabled until every auth instance runs the new
#         KID, then disable the old version.
#
# Never reuse a KID or key material across environments.

SECRET_NAME=${GCP_AUTH_SECRET:-food-flow-auth-keys}
PROJECT_ID=${GCP_PROJECT_ID:-}
WORKLOAD_SERVICE_ACCOUNT=${GCP_WORKLOAD_SERVICE_ACCOUNT:-}

fail() {
    echo "error: $*" >&2
    exit 1
}

for command_name in gcloud openssl jq uuidgen; do
    if ! command -v "${command_name}" >/dev/null 2>&1; then
        fail "missing required command: ${command_name}"
    fi
done

[ -n "${PROJECT_ID}" ] || fail "GCP_PROJECT_ID is not set"

if ! gcloud auth list --filter=status:ACTIVE --format='value(account)' 2>/dev/null | grep -q .; then
    fail "no active gcloud account; run: gcloud auth login"
fi

case "${1:-}" in
    create)
        [ -n "${WORKLOAD_SERVICE_ACCOUNT}" ] || fail "GCP_WORKLOAD_SERVICE_ACCOUNT is required for create"
        if gcloud secrets describe "${SECRET_NAME}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
            fail "secret ${SECRET_NAME} already exists in ${PROJECT_ID}; use 'rotate' to add a new key version"
        fi
        ;;
    rotate)
        if ! gcloud secrets describe "${SECRET_NAME}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
            fail "secret ${SECRET_NAME} does not exist in ${PROJECT_ID}; use 'create' first"
        fi
        ;;
    *)
        echo "usage: $0 {create|rotate}" >&2
        exit 1
        ;;
esac

# Key material lives only in a private temp directory and is shredded on exit.
workdir=$(mktemp -d)
chmod 700 "${workdir}"
cleanup() {
    if command -v shred >/dev/null 2>&1; then
        shred -u "${workdir}/private.pem" "${workdir}/payload.json" 2>/dev/null || true
    else
        rm -fP "${workdir}/private.pem" "${workdir}/payload.json" 2>/dev/null || rm -f "${workdir}/private.pem" "${workdir}/payload.json"
    fi
    rmdir "${workdir}" 2>/dev/null || true
}
trap cleanup EXIT

KID=$(uuidgen)

echo "=> Generating new RSA key pair (kid: ${KID})"
openssl genpkey -algorithm RSA -out "${workdir}/private.pem" -pkeyopt rsa_keygen_bits:2048 2>/dev/null

# jq escapes the PEM newlines into valid JSON; the payload is never echoed.
jq -n --arg key "${KID}" --rawfile pem "${workdir}/private.pem" '{key: $key, pem: $pem}' > "${workdir}/payload.json"

if [ "$1" = "create" ]; then
    echo "=> Creating secret ${SECRET_NAME} in project ${PROJECT_ID}"
    gcloud secrets create "${SECRET_NAME}" --project="${PROJECT_ID}" --replication-policy=automatic
fi

echo "=> Adding a new key version to secret ${SECRET_NAME}"
gcloud secrets versions add "${SECRET_NAME}" --project="${PROJECT_ID}" --data-file="${workdir}/payload.json"

if [ "$1" = "create" ]; then
    echo "=> Granting ${WORKLOAD_SERVICE_ACCOUNT} secretAccessor on ${SECRET_NAME}"
    gcloud secrets add-iam-policy-binding "${SECRET_NAME}" \
        --project="${PROJECT_ID}" \
        --member="serviceAccount:${WORKLOAD_SERVICE_ACCOUNT}" \
        --role="roles/secretmanager.secretAccessor"
fi

cat <<EOF

Done. New key id (KID): ${KID}

Next steps:
  - Deploys must set the active KID to ${KID}:
      staging:    GCP_AUTH_ACTIVE_KID=${KID} ./infra/gcp/deploy-staging.sh
      cloud run:  --update-env-vars AUTH_AUTH_ACTIVE_KID=${KID}
  - For rotation, keep the previous secret version enabled until every auth
    instance runs the new KID, then disable the old version:
      gcloud secrets versions disable <old-version> --secret=${SECRET_NAME} --project=${PROJECT_ID}
EOF
