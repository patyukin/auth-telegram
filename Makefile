.PHONY: start stop rebuild gen-swag

start:
	docker compose up -d

stop:
	docker compose down

rebuild:
	docker compose down -v --remove-orphans
	docker compose up -d --build

gen-swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/auth/main.go -o ./docs/
