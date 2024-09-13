#AppVersion
APP_VERSION=0.1.0

# Name of the binary
BINARY_NAME=service

# Build and binary settings
BUILD_DIR=bin
BUILD_PATH=$(BUILD_DIR)/$(BINARY_NAME)

# Default Go commands
GO=go

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
	-@docker rmi service:$(APP_VERSION)
	@echo "Clean complete."

# Lint Go code (requires golint to be installed)
lint:
	@echo "Linting Go code..."

# Test Go code
unit-test:
	@echo "Running tests..."
	@$(GO) test .
	@echo "Tests complete."

# Test build and pack it into a container
container: lint unit-test build
	mkdir -p docker/bin
	cp bin/service docker/bin
	docker build -t service:$(APP_VERSION) docker

helm: container
	helm package helm/

# Helper make target to rebuild and redeploy on to minikube
refresh-minikube-env: clean helm
	-minikube image rm docker.io/library/service:$(APP_VERSION)
	-minikube image load service:$(APP_VERSION)
	-helm upgrade service service-0.1.0.tgz --set image.pullPolicy='Never'

# Help deploy Prometheus and Grafana
install-monitoring:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	helm install prometheus prometheus-community/prometheus -f monitoring/prometheus-values.yaml
	helm install grafana grafana/grafana -f monitoring/grafana-values.yaml

# Execute the Test.go
send-requests:
	cd test && go build
	./test/test
