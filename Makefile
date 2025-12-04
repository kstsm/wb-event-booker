include .env
export

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Database migrations
migrate-up:
	goose -dir=$(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir=$(MIGRATIONS_DIR) postgres "$(DB_URL)" down

# Start
run:
	go run main.go


