.PHONY: run migrate-up migrate-down migrate-new test test-down lint docker-build up down ps up-build

DOCKER_COMPOSE=docker-compose -f docker/docker-compose.yml
pg_dsn := "postgres://postgres:password@localhost:5432/questions_db?sslmode=disable"

run:
	go run ./cmd/app

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

ps:
	$(DOCKER_COMPOSE) ps

up-build:
	$(DOCKER_COMPOSE) up --force-recreate --remove-orphans -d --build

migrate-new:
	goose -dir migrations create $(name) sql && goose -dir migrations fix

migrate-up:
	goose -dir migrations -env .env up

migrate-down:
	goose -dir migrations -env .env down

test: test-down
	@docker run --rm --name test-questions-pg -p 15432:5432 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=questions_db \
		-d postgres:15-alpine
	@sleep 2
	@goose -timeout 10s -dir migrations postgres "postgres://postgres:password@localhost:15432/questions_db?sslmode=disable" up
	@TEST_POSTGRES_DSN="postgres://postgres:password@localhost:15432/questions_db?sslmode=disable" go test -race -v ./...
	@docker rm -f test-questions-pg

test-down:
	@docker rm -f test-questions-pg 2>/dev/null || true

lint:
	golangci-lint run ./...

docker-build:
	$(DOCKER_COMPOSE) build
