DB_FILE := urls.db
SCHEMA_FILE := schema.sql
GO_BIN := bin/main
MAIN_FILE := main.go

.PHONY: all
all: $(GO_BIN)

$(DB_FILE): $(SCHEMA_FILE)
	@echo "Creating SQLite database and importing schema..."
	sqlite3 $(DB_FILE) < $(SCHEMA_FILE)
	@echo "Database $(DB_FILE) created successfully."

$(GO_BIN): $(DB_FILE) $(MAIN_FILE)
	@echo "Building Go binary..."
	go build -o $(GO_BIN) $(MAIN_FILE)
	@echo "Go application built successfully: $(GO_BIN)"

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(DB_FILE) $(GO_BIN)
	@echo "Clean up complete."