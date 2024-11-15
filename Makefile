watch:
	air -c .air.toml

build:
	export CGO_ENABLED=0
	go build -a -ldflags '-w -extldflags "-static"' -o bin/main cmd/stone-api/main.go

db-up:
	docker compose -f docker/docker-compose.yaml up -d database

db-down:
	docker compose -f docker/docker-compose.yaml down database

db-clean: db-down
	rm -rf docker/mysql/data

db-restart: db-clean db-up

# Preview Local
preview-start: preview-build
	docker compose -f docker/docker-compose.preview.yaml up -d

preview-build:
	docker compose -f docker/docker-compose.preview.yaml build --no-cache

preview-stop:
	docker compose -f docker/docker-compose.preview.yaml down

preview-clean: preview-stop
	docker rmi $(docker images -f "dangling=true" -q)

lint:
	golangci-lint run ./...