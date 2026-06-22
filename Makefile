.PHONY: help bootstrap docker-build docker-build-order docker-build-inventory k8s-load k8s-load-order k8s-load-inventory k8s-apply k8s-destroy nuke

ORDER_IMG ?= order-service:1.16.0
INVENTORY_IMG ?= inventory-service:1.16.0
KIND_CLUSTER_NAME ?= kind

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

bootstrap: k8s-apply docker-build k8s-load ## Full rebuild: terraform apply, then build and load local images

# --- Docker Image Management ---

docker-build: docker-build-order docker-build-inventory ## Build all Docker images

docker-build-order: ## Build the order-service Docker image
	cd order-service && docker build -t $(ORDER_IMG) .

docker-build-inventory: ## Build the inventory-service Docker image
	cd inventory-service && docker build -t $(INVENTORY_IMG) .

# --- Kubernetes Management ---

k8s-load: k8s-load-order k8s-load-inventory ## Load all Docker images into Kind

k8s-load-order: ## Load order-service image into Kind
	kind load docker-image $(ORDER_IMG) --name $(KIND_CLUSTER_NAME)
	# Force pods to restart so they pick up the new image
	kubectl rollout restart deployment/order-service -n default || true

k8s-load-inventory: ## Load inventory-service image into Kind
	kind load docker-image $(INVENTORY_IMG) --name $(KIND_CLUSTER_NAME)
	kubectl rollout restart deployment/inventory-service -n default || true

k8s-apply: ## Run Terraform apply to create/update infrastructure
	cd terraform && terraform init && terraform apply -auto-approve

k8s-destroy: ## Run Terraform destroy to tear down infrastructure
	cd terraform && terraform destroy -auto-approve

# --- Environment Nuke / Clean Slate ---

nuke: ## Nuke the environment completely (destroy cluster and clean terraform state)
	@echo "Nuking the environment..."
	-cd terraform && terraform destroy -auto-approve
	kind delete cluster --name $(KIND_CLUSTER_NAME)
	rm -rf terraform/.terraform terraform/.terraform.lock.hcl terraform/terraform.tfstate terraform/terraform.tfstate.backup
	@echo "Environment nuked successfully. Clean slate achieved."
