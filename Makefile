SHELL = /bin/sh

# запуск докер образов для поднятия локальной инфраструктуры (postgres)
docker-up:
	PWD=$(PWD) user=$(id -u) group=$(id -g) docker-compose --file build/docker-compose.yaml up

dependencies:
	go build \
    		-ldflags="-s -w" \
            -tags='no_sqlite no_clickhouse no_mssql no_mysql' \
            -o goose ./bin/goose

migrations:
	@goose --dir=./migrations/postgres postgres 'postgres://postgres:pg-pass@localhost:6432/pg-db?sslmode=disable' up

# запуск приложения
run:
	go run main.go

docker-down:
	docker-compose --file build/docker-compose.yaml down