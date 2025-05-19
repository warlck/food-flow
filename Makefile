# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Define dependencies

KIND_CLUSTER    := food-flow-cluster
KIND            := kindest/node:v1.32.2
GOLANG          := golang:1.24
ALPINE          := alpine:3.21
POSTGRES        := postgres:17.4
GRAFANA         := grafana/grafana:11.6.0
PROMETHEUS      := prom/prometheus:v3.2.0
TEMPO           := grafana/tempo:2.7.0
LOKI            := grafana/loki:3.4.0
PROMTAIL        := grafana/promtail:3.4.0

NAMESPACE       := sales-system
SALES_APP       := sales-api
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/food-flow
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"
# ==============================================================================
# Building containers

build: sales 

sales:
	docker build \
		-f infra/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.


# ==============================================================================
# Running the service
run:
	go run ./app/services/sales-api/main.go | go run ./app/tooling/logfmt/main.go

help:
	go run ./app/services/sales-api/main.go --help

version:
	go run ./app/services/sales-api/main.go --version



# ==============================================================================
# Debugging


debug-statsviz:
	open http://localhost:3010/debug/statsviz



# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner


dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)


dev-apply:
	kustomize build infra/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready


dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-run: build dev-up dev-load dev-apply

dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

# ------------------------------------------------------------------------------

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SALES_APP)

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(SALES_APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)


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




