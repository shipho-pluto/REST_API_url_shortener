include .env
export

migrate-up:
	go run ./cmd/migrator/main.go -command="up" -dir=$(MIGRATION_PATH) -config_path=$(CONFIG_PATH)
migrate-down:
	go run ./cmd/migrator/main.go -command="down" -dir=$(MIGRATION_PATH) -config_path=$(CONFIG_PATH)
migrate-refresh:
	go run ./cmd/migrator/main.go -command="refresh" -dir=$(MIGRATION_PATH) -config_path=$(CONFIG_PATH)
