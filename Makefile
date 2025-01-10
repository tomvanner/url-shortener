include .env

COMPOSE_FILE=docker-compose.yaml
GO_BIN := bin/app
MAIN_FILE := app.go

default: db-up

.PHONY: db
db: $(DB_PATH)
$(DB_PATH): $(SCHEMA_PATH)
	mkdir -p $(dir $@)
	if [ ! -f $@ ]; then \
		sqlite3 $@ < $<; \
		echo "Database created at $@"; \
	else \
		echo "Database already exists at $@, skipping creation"; \
	fi

.PHONY: up
up:
	@echo "Building and starting docker containers"
	docker-compose -f $(COMPOSE_FILE) up --build



.PHONY: down
down:
	@echo "Stopping and removing docker containers"
	docker-compose -f $(COMPOSE_FILE) down

.PHONY: app
app:
	$(GO_BIN): $(MAIN_FILE)
		@echo "Building Go binary..."
		go build -o $(GO_BIN) $(MAIN_FILE)
		@echo "Go application built successfully: $(GO_BIN)"

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(DB_PATH)
	@echo "Clean up complete"

.PHONY: db-up
db-up: db up