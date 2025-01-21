COMPOSE_FILE=docker-compose.yml

APP_SERVICE=my-app
DB_SERVICE=my-postgres

.PHONY: build up start stop down logs clean

build:
	docker-compose -f $(COMPOSE_FILE) build

up:
	docker-compose -f $(COMPOSE_FILE) up -d

start:
	docker-compose -f $(COMPOSE_FILE) start

stop:
	docker-compose -f $(COMPOSE_FILE) stop

down:
	docker-compose -f $(COMPOSE_FILE) down

logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

clean: down
	docker-compose -f $(COMPOSE_FILE) down -v