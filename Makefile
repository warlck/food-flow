# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem
# ==============================================================================

.PHONY: \
	build sales auth frontend gcp-staging-deploy \
	run help version \
	curl-live curl-ready curl-test-error load-test curl-auth curl-create-user \
	admin-genkey pgcli \
	debug-statsviz metrics \
	dev-up dev-down dev-load-db dev-load dev-apply dev-restart dev-run dev-update dev-update-apply dev-stripe-secrets \
	dev-logs dev-logs-auth dev-logs-frontend dev-describe-deployment dev-describe-sales dev-describe-auth dev-describe-frontend dev-logs-db

# Define dependencies

KIND_CLUSTER    := food-flow-cluster
KIND            := kindest/node:v1.34.0
GOLANG          := golang:1.25
ALPINE          := alpine:3.21
POSTGRES        := postgres:18
GRAFANA         := grafana/grafana:11.6.0
PROMETHEUS      := prom/prometheus:v3.2.0
TEMPO           := grafana/tempo:2.7.0
LOKI            := grafana/loki:3.4.0
PROMTAIL        := grafana/promtail:3.4.0

NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/food-flow
VERSION         := 0.0.5
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)
FRONTEND_IMAGE  := $(BASE_IMAGE_NAME)/frontend:$(VERSION)
HELM_CHART      := infra/helm/food-flow
HELM_RELEASE    := food-flow
HELM_DEV_VALUES := $(HELM_CHART)/values-kind.yaml
DATABASE_SECRET := $(HELM_RELEASE)-database-credentials
STRIPE_SECRET   := $(HELM_RELEASE)-stripe-secrets

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"
# ==============================================================================
# Building containers

build: sales auth frontend

sales:
	docker build \
		-f infra/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

auth:
	docker build \
		-f infra/docker/dockerfile.auth \
		-t $(AUTH_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

frontend:
	docker build \
		-f infra/docker/dockerfile.frontend \
		-t $(FRONTEND_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg VITE_API_URL="" \
		--build-arg VITE_STRIPE_PUBLISHABLE_KEY="$(VITE_STRIPE_PUBLISHABLE_KEY)" \
		.

gcp-staging-deploy:
	./infra/gcp/deploy-staging.sh


# ==============================================================================
# Running the service
run:
	go run ./api/services/sales-api/main.go | go run ./api/tooling/logfmt/main.go

help:
	go run ./api/services/sales-api/main.go --help

version:
	go run ./app/services/sales-api/main.go --version

curl-live:
	curl -il -X GET http://localhost:3000/v1/liveness

curl-ready:
	curl -il -X GET http://localhost:3000/v1/readiness

curl-test-error:
	curl -il -X GET http://localhost:3000/v1/testerror

load-test:
	hey -m GET -c 100 -n 100000 "http://localhost:3000/v1/testerror"

admin-genkey:
	go run ./api/tooling/sales-admin/main.go

curl-auth:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/testauth"

curl-create-user:
	curl -il -X POST -H 'Content-Type: application/json' -d '{"name":"foo","email":"foo@bar.com","roles":["ADMIN"],"department":"IT","password":"123","passwordConfirm":"123"}' http://localhost:3000/v1/users


pgcli:
	pgcli postgresql://postgres:postgres@localhost






# ==============================================================================
# Metrics and Tracing

debug-statsviz:
	open http://localhost:3010/debug/statsviz


metrics:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"



# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner
		
# 	docker save $(POSTGRES) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & 
# 	wait;





dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ------------------------------------------------------------------------------

dev-load-db:
	@echo "Note: Skipping image pre-load due to Kind/containerd issue. Image will be pulled directly by Kubernetes."

dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)
	kind load docker-image $(AUTH_IMAGE) --name $(KIND_CLUSTER)
	kind load docker-image $(FRONTEND_IMAGE) --name $(KIND_CLUSTER)


dev-stripe-secrets:
	@if [ -z "$$SALES_STRIPE_SECRET_KEY" ]; then echo "SALES_STRIPE_SECRET_KEY is not set"; exit 1; fi
	@if [ -n "$$SALES_STRIPE_WEBHOOK_SECRET" ]; then \
		kubectl create secret generic $(STRIPE_SECRET) --namespace=$(NAMESPACE) \
			--from-literal=stripe_secret_key="$$SALES_STRIPE_SECRET_KEY" \
			--from-literal=stripe_webhook_secret="$$SALES_STRIPE_WEBHOOK_SECRET" \
			--dry-run=client -o yaml | kubectl apply -f - ; \
	else \
		echo "SALES_STRIPE_WEBHOOK_SECRET is not set; creating secret without webhook secret" ; \
		kubectl create secret generic $(STRIPE_SECRET) --namespace=$(NAMESPACE) \
			--from-literal=stripe_secret_key="$$SALES_STRIPE_SECRET_KEY" \
			--dry-run=client -o yaml | kubectl apply -f - ; \
	fi


dev-apply:
	@if [ -z "$$POSTGRES_PASSWORD" ]; then echo "POSTGRES_PASSWORD is not set"; exit 1; fi
	kubectl create namespace $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	kubectl create secret generic $(DATABASE_SECRET) --namespace=$(NAMESPACE) \
		--from-literal=password="$$POSTGRES_PASSWORD" --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade --install $(HELM_RELEASE) $(HELM_CHART) \
		--namespace=$(NAMESPACE) \
		--values=$(HELM_DEV_VALUES) \
		--wait \
		--timeout=300s
	
	# If image tags don't change (e.g. localhost/food-flow/frontend:0.0.5), Kubernetes won't update pods
	# just because you rebuilt and re-loaded the image into kind. Force a restart to pick up the new image.
	kubectl rollout restart deployment $(HELM_RELEASE)-$(SALES_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(HELM_RELEASE)-$(AUTH_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(HELM_RELEASE)-frontend --namespace=$(NAMESPACE)
	
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s deployment/$(HELM_RELEASE)-$(SALES_APP)
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s deployment/$(HELM_RELEASE)-$(AUTH_APP)
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s deployment/$(HELM_RELEASE)-frontend


dev-restart:
	kubectl rollout restart deployment $(HELM_RELEASE)-$(SALES_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(HELM_RELEASE)-$(AUTH_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(HELM_RELEASE)-frontend --namespace=$(NAMESPACE)
 
dev-run: build dev-up dev-load dev-apply

dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

# ------------------------------------------------------------------------------

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run api/tooling/logfmt/main.go -service=$(SALES_APP)

dev-logs-auth:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(AUTH_APP) --all-containers=true -f --tail=100 | go run api/tooling/logfmt/main.go -service=$(AUTH_APP)

dev-logs-frontend:
	kubectl logs --namespace=$(NAMESPACE) -l app=frontend --all-containers=true -f --tail=100

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(HELM_RELEASE)-$(SALES_APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)

dev-describe-auth:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(AUTH_APP)

dev-describe-frontend:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=frontend

dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

# ------------------------------------------------------------------------------
dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) -f --tail=100 -c init-migrate-seed


# ------------------------------------------------------------------------------
# Status

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces




# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

test-race:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	@docker inspect servicetest >/dev/null 2>&1 || docker run -d --name servicetest -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:17.5 -c log_statement=all
	@until docker exec servicetest pg_isready -U postgres >/dev/null 2>&1; do sleep 1; done
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck ./...

vuln-check:
	-govulncheck ./...

test: test-only lint vuln-check

test-race: test-race lint vuln-check

# ==============================================================================
# Hitting endpoints

token:
	curl -i \
	--user "admin@example.com:gophers" http://localhost:6000/v1/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1


users:
	curl -i \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"
	
test-auth:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:6000/v1/auth/authenticate"
