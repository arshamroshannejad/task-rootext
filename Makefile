up:
	docker compose up --build -d

down:
	docker compose down

log-server:
	docker compose logs -f server

log-postgres:
	docker compose logs -f postgres

log-redis:
	docker compose logs -f redis
