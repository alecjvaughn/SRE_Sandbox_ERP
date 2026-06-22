.PHONY: help bootstrap docker-build docker-build-order docker-build-inventory k8s-load k8s-load-order k8s-load-inventory k8s-apply k8s-destroy nuke

ORDER_IMG ?= order-service:1.16.0
INVENTORY_IMG ?= inventory-service:1.16.0
PROXY_IMG ?= alecjvaughn/edge-proxy:latest
KIND_CLUSTER_NAME ?= kind

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

bootstrap: k8s-create k8s-apply docker-build k8s-load ## Full rebuild: create cluster, terraform apply, then build and load local images

# --- Docker Image Management ---

docker-build: docker-build-order docker-build-inventory docker-build-proxy ## Build all Docker images

docker-build-order: ## Build the order-service Docker image
	cd order-service && docker build -t $(ORDER_IMG) .

docker-build-inventory: ## Build the inventory-service Docker image
	cd inventory-service && docker build -t $(INVENTORY_IMG) .

docker-build-proxy: ## Build the edge-proxy Docker image
	cd edge-proxy && docker build -t $(PROXY_IMG) .

# --- Kubernetes Management ---
k8s-create: ## Create the Kind cluster
	kind create cluster --name $(KIND_CLUSTER_NAME) || true

k8s-load: k8s-load-order k8s-load-inventory ## Load all Docker images into Kind

k8s-load-order: ## Load order-service image into Kind
	kind load docker-image $(ORDER_IMG) --name $(KIND_CLUSTER_NAME)
	# Force pods to restart so they pick up the new image
	kubectl rollout restart deployment/order-service -n default || true

k8s-load-inventory: ## Load inventory-service image into Kind
	kind load docker-image $(INVENTORY_IMG) --name $(KIND_CLUSTER_NAME)
	kubectl rollout restart deployment/inventory-service -n default || true

# --- Azure Infrastructure Management ---
azure-apply: ## Run Terraform apply to create/update Azure infrastructure
	cd terraform && terraform init -upgrade && terraform apply -auto-approve
	# Get AKS credentials so kubectl works
	az aks get-credentials --resource-group rg-sre-sandbox --name aks-sre-sandbox --overwrite-existing
	kubectl apply -f argocd/project.yaml
	kubectl apply -f argocd/root-app.yaml

azure-destroy: ## Run Terraform destroy to tear down Azure infrastructure
	cd terraform && terraform destroy -auto-approve

# --- Environment Nuke / Clean Slate ---

nuke: ## Nuke the Azure environment completely (destroy cluster and clean terraform state)
	@echo "Nuking the environment..."
	-cd terraform && terraform destroy -auto-approve
	# Fallback if terraform state is corrupted
	-az group delete --name rg-sre-sandbox --yes --no-wait
	rm -rf terraform/.terraform terraform/.terraform.lock.hcl terraform/terraform.tfstate terraform/terraform.tfstate.backup
	@echo "Environment nuked successfully. Clean slate achieved."
