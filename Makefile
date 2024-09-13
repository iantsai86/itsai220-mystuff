# Name of the binary
BINARY_NAME=service

# Build and binary settings
BUILD_DIR=bin
BUILD_PATH=$(BUILD_DIR)/$(BINARY_NAME)

# Default Go commands
GO=go
GOREBUILD=go build -o $(BUILD_PATH)

# Build the Go binary
build:
	@echo "Building the binary..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm $(GO) build -o $(BUILD_PATH) ./main.go
	@echo "Build complete: $(BUILD_PATH)"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	-@rm -rf $(BUILD_DIR)
	-@rm -rf docker/$(BUILD_DIR)
	-@rm service-*.tgz
	-@docker rmi service:latest
	@echo "Clean complete."

# Lint Go code (requires golint to be installed)
lint:
	@echo "Linting Go code..."

# Test Go code
unit-test:
	@echo "Running tests..."
	@$(GO) test .
	@echo "Tests complete."

container: lint unit-test build
	mkdir -p docker/bin
	cp bin/service docker/bin
	docker build -t service:latest docker

helm: container
	helm package helm/

refresh-minikube-env: clean helm
	-helm uninstall service 
	sleep 3
	-minikube image rm docker.io/library/service:latest
	-minikube image load service:latest
	-helm install service service-0.1.0.tgz --set image.pullPolicy='Never'

install-monitoring:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	helm install prometheus prometheus-community/prometheus -f monitoring/prometheus-values.yaml
	helm install grafana grafana/grafana -f monitoring/grafana-values.yaml
	