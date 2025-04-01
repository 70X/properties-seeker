.PHONY: dev test

test:
	docker compose up --build test-db

dev:
	docker compose up --build dev-db
