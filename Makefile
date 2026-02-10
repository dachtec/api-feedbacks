.PHONY: build run dev test lint docker-up docker-down docker-dev clean

# Build the application binary inside Docker
build:
	docker compose build app

# Run production containers
run: docker-up

# Run development containers with hot-reload
dev:
	docker compose --profile dev up --build

# Run tests inside a Docker container
test:
	docker compose run --rm --no-deps app-dev sh -c "go test ./internal/... -v -count=1"

# Run tests with coverage inside Docker
test-cover:
	docker compose run --rm --no-deps app-dev sh -c "go test ./internal/... -v -coverprofile=coverage.out && go tool cover -func=coverage.out"

# Run linter inside Docker
lint:
	docker compose run --rm --no-deps app-dev sh -c "go vet ./..."

# Start production environment
docker-up:
	docker compose up -d --build

# Stop all containers
docker-down:
	docker compose down

# Stop all containers and remove volumes
docker-clean:
	docker compose down -v

# Display logs
logs:
	docker compose logs -f app

# Seed sample data
seed:
	@echo "Seeding sample data..."
	@bash scripts/seed.sh

clean:
	docker compose down -v --rmi local
	rm -rf tmp/
