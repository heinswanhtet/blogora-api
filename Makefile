include .env
export

MIGRATIONS_DIR = ./migrations
DB_DRIVER = mysql
DB_URL = $(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&loc=UTC

build:
	@go build -o ./bin/blogora-api ./cmd/main.go

run: build
	@./bin/blogora-api

test:
	@go test -v ./tests

tidy:
	@go mod tidy

vendor: 
	@go mod tidy; 
	go mod vendor

migration:
	@echo "❌ Please run like: make migration-add_your_name"

migration-%:
	@mkdir -p $(MIGRATIONS_DIR)
	@filename="$$(date +%Y%m%d%H%M%S)_$*.sql"; \
	touch "$(MIGRATIONS_DIR)/$$filename"; \
	echo "-- +goose Up\n\n-- +goose Down" > "$(MIGRATIONS_DIR)/$$filename"; \
	echo "✅ Created $(MIGRATIONS_DIR)/$$filename"

# goose
up:
	@goose -dir $(MIGRATIONS_DIR) "$(DB_DRIVER)" "$(DB_URL)" up

down:
	@goose -dir $(MIGRATIONS_DIR) "$(DB_DRIVER)" "$(DB_URL)" down

status:
	@goose -dir $(MIGRATIONS_DIR) "$(DB_DRIVER)" "$(DB_URL)" status

create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql

# air
air:
	@$$HOME/go/bin/air