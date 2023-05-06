SHELL = /bin/sh

# параметры для конфига прописаны в main.go (http, db параметры)
# миграции для БД находятся в migrations/table.sql
# тестовые данные для БД находятся в migrations/test_data.sql

# запуск докер образов для поднятия локальной инфраструктуры (postgres)
docker_up:
	PWD=$(PWD) user=$(id -u) group=$(id -g) docker-compose --file build/docker-compose.yaml up

# запуск приложения
run:
	go run main.go