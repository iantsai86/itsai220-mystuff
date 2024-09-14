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
	-minikube image load $(BINARY_NAME):$(APP_VERSION)
	@if helm status $(BINARY_NAME) > /dev/null 2>&1; then \
		helm upgrade $(BINARY_NAME) $(BINARY_NAME)-$(APP_VERSION).tgz --set image.pullPolicy='Never' ; \
	else \
		echo "Helm installing $(BINARY_NAME)" ; \
		helm install $(BINARY_NAME) $(BINARY_NAME)-$(APP_VERSION).tgz --set image.pullPolicy='Never' ; \
	fi

# Deploy Kube Prometheus Stack 
install-monitoring:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	helm install prometheus prometheus-community/kube-prometheus-stack

# Execute the Test.go
send-requests:
	cd test && go build
	./test/test
