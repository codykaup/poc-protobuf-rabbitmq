help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## Build the protobuf message
	protoc --go_out=./generated --go_opt=paths=source_relative message.proto

run: ## Run the project
	go run main.go
