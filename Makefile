SHELL = /bin/sh

# запуск докер образов для поднятия локальной инфраструктуры (postgres)
docker_up:
	PWD=$(PWD) user=$(id -u) group=$(id -g) docker-compose --file build/docker-compose.yaml up

# запуск приложения
run:
	go run main.go