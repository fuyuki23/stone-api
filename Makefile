watch:
	air -c .air.toml

build:
	export CGO_ENABLED=0
	go build -a -ldflags '-w -extldflags "-static"' -o bin/main cmd/stone-api/main.go

db-up:
	docker compose -f docker/docker-compose.yaml up -d

db-down:
	docker compose -f docker/docker-compose.yaml down

db-clean: db-down
	rm -rf docker/mysql/data

db-restart: db-clean db-up
