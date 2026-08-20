# GCP deployment and auth runbook

This directory contains the Cloud Run deploy scripts and the auth signing-key
bootstrap. Staging and production are fully separate: own projects, Artifact
Registry repos, secrets, and workload service accounts. Nothing is shared.

## Environments

| | Staging | Production |
|---|---|---|
| Deploy script | `deploy-staging.sh` | `deploy-production.sh` |
| Auth key secret | `food-flow-auth-keys` (staging project) | `food-flow-auth-keys` (production project) |
| Admin frontend | public URL, login page only | gated by IAP (Google account) **and** the login page |
| Auth service | min-instances 0, CPU throttled | min-instances 1, no throttling (sales delegates authentication to it on every call) |
| CORS on auth | `*` default | pinned to the admin origin by the deploy script |

## Bootstrap the auth signing key (once per environment)

```sh
export GCP_PROJECT_ID=<project>
infra/gcp/bootstrap-auth-secret.sh create
```

The script generates an RSA keypair locally, stores `{"key":"<kid>","pem":"..."}`
in Secret Manager as `food-flow-auth-keys`, grants the workload service account
`secretAccessor`, and prints the KID. The private key never leaves the machine
except into the secret payload; the temp files are shredded.

Export the printed KID for deploys:

```sh
export GCP_AUTH_ACTIVE_KID=<kid from bootstrap>
```

The deploy scripts fail fast without it. The auth service also fails fast at
startup if `AUTH_AUTH_ACTIVE_KID` is not present in the keystore.

## Deploy

```sh
# staging
GCP_PROJECT_ID=<staging-project> GCP_AUTH_ACTIVE_KID=<kid> infra/gcp/deploy-staging.sh

# production
GCP_PROJECT_ID=<prod-project> GCP_AUTH_ACTIVE_KID=<kid> infra/gcp/deploy-production.sh
```

Both scripts build and push images, run migrations, deploy auth (with
`--set-secrets AUTH_AUTH_KEYS_ENV_VAR=food-flow-auth-keys:latest`), sales, and
the storefront, then deploy the admin frontend and verify health endpoints.

Production additionally pins `AUTH_WEB_CORS_ALLOWED_ORIGINS` to the admin
origin and enables IAP on the admin frontend. One-time IAP setup (console):
configure the OAuth consent screen, then grant each operator the
**IAP-secured Web App User** role on the admin service. Programmatic checks
against the gated admin frontend need an ID token minted for the IAP OAuth
client audience (`gcloud auth print-identity-token --audiences=<IAP_CLIENT_ID>`).

## Key rotation

Rotate when a key may be exposed, when an operator with secret access leaves,
or on a periodic schedule:

```sh
infra/gcp/bootstrap-auth-secret.sh rotate
```

This adds a **new** secret version with a fresh KID and prints it. Redeploy
with the new KID:

```sh
GCP_AUTH_ACTIVE_KID=<new kid> infra/gcp/deploy-production.sh
```

Cutover impact: the auth service restarts with the new key and immediately
starts signing with it; tokens signed by the old key stop verifying (the old
key is no longer in the keystore), so every admin is forced to log in again.
That is deliberate — rotation is also the emergency "kill all sessions"
button. Zero-downtime rotation (keystore holds old+new keys during a grace
window) is a documented future enhancement.

## User provisioning and lifecycle

Never run `seed.sql` outside kind — the well-known seed credentials must not
exist in staging or production.

**Create an admin user** (one-off Cloud Run job, same pattern as migrations):

```sh
gcloud run jobs deploy production-db-useradd \
    --project <project> --region <region> \
    --image <sales image> --command ./admin \
    --args "useradd,--name,<Name>,--email,<email>,--password,<password>" \
    --service-account <workload SA> \
    --network <vpc> --subnet <subnet> --vpc-egress private-ranges-only \
    --set-secrets "SALES_DB_PASSWORD=food-flow-db-password:latest" \
    --execute-now --wait
```

Passwords must satisfy the policy: 12-64 printable ASCII characters.

**Reset a password:** `admin useradd` is create-only. Use the sales API with
an admin token (log in to Restaurant Studio and take the issued token). The
sales service is internal-ingress, so tunnel through Cloud Run first:

```sh
gcloud beta run services proxy production-sales --project <project> --region <region> --port 3000

curl -X PUT "http://localhost:3000/v1/users/<user-id>" \
    -H "Authorization: Bearer <admin token>" -H "Content-Type: application/json" \
    -d '{"password":"<new password>","passwordConfirm":"<new password>"}'
```

**Disable / offboard:** set `enabled=false` via `PUT /v1/users/<user-id>`.
The auth service checks the database on every authentication, so a disabled
user's existing tokens stop working on their next request.

**Emergency response:** disable the user (above) and rotate the signing key
(above) to invalidate every outstanding token immediately.

## Alerting

The auth service emits a structured log marker on every error-level event:

```
msg="SEND ALERT" alert=true errorMessage=...
```

Wire a GCP log-based metric on `jsonPayload.alert="true"` filtered to the
auth service, and alert on any sustained rate. Complement with:

- 429 responses from the admin nginx (`limit_req` rejections) — brute-force
  signal, visible in the admin frontend access logs.
- Auth service unavailable / 5xx — because sales delegates authentication to
  auth, an auth outage is an admin-API outage (this is why production runs
  auth with min-instances 1).

## Local fallback path

If `AUTH_AUTH_KEYS_ENV_VAR` is unset, the auth service loads PEMs from
`infra/keys/` (KID = filename without `.pem`). This exists for offline
`go run` development only. The folder is gitignored (`infra/keys/*.pem`) and
the auth Docker image no longer ships any keys. Never commit PEMs.

## Committed-key history note

The original development key (`54bb2165-…pem`) was committed to this
repository before per-environment keys existed. It has been removed from the
tree and the image, and rotation neutralizes it, but it remains readable in
git history. Rewriting history (force-push, breaking all clones/forks/CI
caches) is deliberately out of scope; if a purge is ever required, follow the
filter-repo procedure in the spec's Appendix A (`docs/specs/login-gcp-secrets.md`).
