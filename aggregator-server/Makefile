include .env

MIGRATE = migrate -database postgres://root:secret@localhost:5432/aggregator_db?sslmode=disable -path migrations

stop_containers:
	@echo "Stopping other docker containers"
	if [ $$(docker ps -q)]; then \
		echo "found and stopped containers" \
		docker stop $$(docker ps -q); \
	else \
		echo "no containers running..."; \
	fi

create_container:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=${USER} -e POSTGRES_PASSWORD=${PASSWORD} -d postgres:12-alpine

create_db:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --usernam=${USER} --owner=${USER} ${DB_NAME}

start_container:
	docker start ${DB_DOCKER_CONTAINER}

create_migrations:
	${MIGRATE} create -ext sql -dir migrations/ -seq add_chapters

migrate_up:
	${MIGRATE} up

migrate_down:
	${MIGRATE} down

build:
	if [ -f "${BINARY}" ]; then \
		rm ${BINARY}; \
		echo "Deleted ${BINARY}"; \
	fi	
	@echo "Building binary..."
	go build -o ${BINARY} cmd/server/*.go

run: build
	./${BINARY}	
# @echo "api started..."

stop:
	@echo "stopping server.."
	@-pkill -SIGTERM -f "./${BINARY}"
	@echo "server stopped..."