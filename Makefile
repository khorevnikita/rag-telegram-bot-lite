MIGRATIONS_PATH=./cmd/migrations/main.go
SEED_PATH=./cmd/seeder/main.go
NETWORK=$(shell grep '^NETWORK_NAME=' .env | cut -d '=' -f 2)

# Команда для запуска миграций (up)
migrate-up:
	@echo "Using network: $(NETWORK)"
	docker run --rm \
		--network $(NETWORK) \
		-v $(PWD)/bot:/app \
		-w /app \
		-v $(PWD)/.env:/app/.env \
		golang:1.23 \
		go run $(MIGRATIONS_PATH) -command=up

# Команда для отката миграций (down)
migrate-down:
	@echo "Using network: $(NETWORK)"
	docker run --rm \
		--network $(NETWORK) \
		-v $(PWD)/bot:/app \
		-w /app \
		-v $(PWD)/.env:/app/.env \
		golang:1.23 \
		go run $(MIGRATIONS_PATH) -command=down

migration:
	@echo "Using network: $(NETWORK)"
	docker run --rm -v $(PWD)/bot:/app -w /app golang:1.23 \
	sh -c "go install github.com/pressly/goose/v3/cmd/goose@latest && /go/bin/goose -dir ./database/migrations create \"$$name\" sql"

seed-form:
	@echo "Using network: $(NETWORK)"
	docker run --rm \
		--network $(NETWORK) \
		-v $(PWD)/bot:/app \
		-w /app \
		-v $(PWD)/.env:/app/.env \
		golang:1.23 \
		go run $(SEED_PATH) questions

seed-contexts:
	@echo "Using network: $(NETWORK)"
	docker run --rm \
		--network $(NETWORK) \
		-v $(PWD)/bot:/app \
		-w /app \
		-v $(PWD)/.env:/app/.env \
		golang:1.23 \
		go run $(SEED_PATH) contexts

seed-payment-providers:
	@echo "Using network: $(NETWORK)"
	docker run --rm \
		--network $(NETWORK) \
		-v $(PWD)/bot:/app \
		-w /app \
		-v $(PWD)/.env:/app/.env \
		golang:1.23 \
		go run $(SEED_PATH) payment_providers

lint:
	cd bot && \
	golangci-lint run --config .golangci.yml ./...