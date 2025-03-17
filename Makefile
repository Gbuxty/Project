CONFIG_PATH=config/local.yaml

AUTH_DB_DSN=postgres://postgres:postgres_1234@localhost:5432/Authentication?sslmode=disable
AUTH_MIGRATIONS_DIR=./AuthService/migrations/postgres

FEED_DB_DSN=postgres://feeduser:feedpassword@localhost:5433/feedservice?sslmode=disable
FEED_MIGRATIONS_DIR=./FeedService/migrations/postgres

gen-proto:
	protoc --go_out=./proto/ --go-grpc_out=./proto/ ./proto/authentication_feed.proto

run:
	go run ./cmd/main.go -c=$(CONFIG_PATH)


migrate-create-auth-%:
	goose -dir $(AUTH_MIGRATIONS_DIR) create $(subst migrate-create-auth-,,$@) sql

migrate-up-auth:
	goose -dir $(AUTH_MIGRATIONS_DIR) postgres "$(AUTH_DB_DSN)" up

migrate-down-auth:
	goose -dir $(AUTH_MIGRATIONS_DIR) postgres "$(AUTH_DB_DSN)" down

migrate-reset-auth:
	goose -dir $(AUTH_MIGRATIONS_DIR) postgres "$(AUTH_DB_DSN)" reset

migrate-status-auth:
	goose -dir $(AUTH_MIGRATIONS_DIR) postgres "$(AUTH_DB_DSN)" status


migrate-create-feed-%:
	goose -dir $(FEED_MIGRATIONS_DIR) create $(subst migrate-create-feed-,,$@) sql

migrate-up-feed:
	goose -dir $(FEED_MIGRATIONS_DIR) postgres "$(FEED_DB_DSN)" up

migrate-down-feed:
	goose -dir $(FEED_MIGRATIONS_DIR) postgres "$(FEED_DB_DSN)" down

migrate-reset-feed:
	goose -dir $(FEED_MIGRATIONS_DIR) postgres "$(FEED_DB_DSN)" reset

migrate-status-feed:
	goose -dir $(FEED_MIGRATIONS_DIR) postgres "$(FEED_DB_DSN)" status


docker-build:
	docker build -t your-app-name .

docker-run:
	docker run -p 9090:9090 your-app-name

docker-compose-up:
	docker-compose up --build

help:
	@echo "Доступные команды:"
	@echo "  gen-proto          - Генерация protobuf"
	@echo "  run                - Запуск приложения"
	@echo "  migrate-create-auth-% - Создать новую миграцию для AuthService (замените % на имя миграции)"
	@echo "  migrate-up-auth    - Применить миграции для AuthService"
	@echo "  migrate-down-auth  - Откатить последнюю миграцию для AuthService"
	@echo "  migrate-reset-auth - Откатить все миграции для AuthService"
	@echo "  migrate-status-auth - Показать статус миграций для AuthService"
	@echo "  migrate-create-feed-% - Создать новую миграцию для FeedService (замените % на имя миграции)"
	@echo "  migrate-up-feed    - Применить миграции для FeedService"
	@echo "  migrate-down-feed  - Откатить последнюю миграцию для FeedService"
	@echo "  migrate-reset-feed - Откатить все миграции для FeedService"
	@echo "  migrate-status-feed - Показать статус миграций для FeedService"
	@echo "  docker-build       - Собрать Docker-образ"
	@echo "  docker-run         - Запустить Docker-контейнер"
	@echo "  docker-compose-up  - Запустить docker-compose"
	@echo "  help               - Показать эту справку"