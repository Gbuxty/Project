DB_DSN=postgres://postgres:postgres_1234@localhost:5432/Authentication?sslmode=disable
MIGRATIONS_DIR=./migrations/postgres


gen-proto:
	protoc --go_out=./proto/ --go-grpc_out=./proto/ ./proto/authentication.proto

run:
	go run ./cmd/main.go

migrate-create-%:
	goose -dir $(MIGRATIONS_DIR) create $(subst migrate_create_,,$@) sql

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up


migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" reset

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

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
	@echo "  migrate_create_%   - Создать новую миграцию (замените % на имя миграции)"
	@echo "  migrate_up         - Применить миграции"
	@echo "  migrate_down       - Откатить последнюю миграцию"
	@echo "  migrate_reset      - Откатить все миграции"
	@echo "  migrate_status     - Показать статус миграций"
	@echo "  docker-build       - Собрать Docker-образ"
	@echo "  docker-run         - Запустить Docker-контейнер"
	@echo "  docker-compose-up  - Запустить docker-compose"
	@echo "  help               - Показать эту справку"