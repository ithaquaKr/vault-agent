# Makefile for vault-agent
# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

##@ General
# Targets commented with ## will be visible in "make help" info.
# Comments marked with ##@ will be used as categories for a group of targets.

.PHONY: help
default: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: compose-up compose-down
compose-up: ## Start docker compose development environment
	docker compose -f ./tests/compose.yml up -d

compose-down: ## Destroy docker compose development environment
	docker compose -f ./tests/compose.yml down -v

PORT_FORWARDING=8200 # This port is hardcode in Vault Helm values file, if you want to change, change values.yaml first.

.PHONY: kind-up kind-down
kind-up: ## Start kubernetes development environment using kind
	@echo "Creating Kind cluster..."
	kind create cluster --config=./tests/kind/config.yaml
kind-down: ## Destroy kubernetes development environment.
	@echo "Destroying Kind cluster..."
	# Hardcode name ??
	kind delete cluster -n kind-cluster

.PHONY: vault-up vault-down port-forward
vault-up: ## Install Vault Cluster using Helm to current Kubernetes cluster context.
	@echo "Installing Vault using Helm"
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm upgrade --install vault hashicorp/vault -n vault --create-namespace --values ./tests/kind/vault-helm/values.yaml
	@echo "Vault development environment in Kubernetes installed."

vault-down: ## Destroy Vault Cluster
	kubectl delete namespace vault
	@echo "Kubernetes Vault development environment destroyed."

port-forward: ## Port forwarding for Vault that install in Kubernetes
	@echo "Port-forwarding Vault port ${PORT_FORWARDING} to localhost:${PORT_FORWARDING}"
	kubectl port-forward -n vault services/vault ${PORT_FORWARDING}:${PORT_FORWARDING}

##@ Build
BINARY_NAME=vault-agent
VERSION?=0.1.0
GO_VERSION=$(shell go version | awk '{print $$3}')
COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"

.PHONY: build clean
build: ## Build binary
	go build -race $(LDFLAGS) -o ./bin/$(BINARY_NAME) ./cmd/cli

clean: ## Clean binary file
	rm -f ./bin/$(BINARY_NAME)
	go clean -cache

##@ Test
.PHONY: test
test: ## Run tests
	go test -race -v ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	go test -race -v -tags=integration ./...

##@ Dependencies
.PHONY: deps
deps: ## Install dependencies
	go mod tidy
	go mod verify
