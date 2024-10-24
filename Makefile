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
	docker run --name ${DB_DOCKER_CONTAINER_NAME} -p 5432:5432 -e POSTGRES_USER=${USER} -e POSTGRES_PASSWORD=${PASSWORD} -d postgres:12-alpine

create_db:
	docker exec -it ${DB_DOCKER_CONTAINER_NAME} createdb --username=${USER} --owner=${USER} ${DB_NAME}

start_container:
	docker start ${DB_DOCKER_CONTAINER_NAME}

create_migrations:
	${MIGRATE} create -ext sql -dir migrations/ -seq add_vector_extension_and_tables

migrate_up:
	${MIGRATE} up

migrate_down:
	${MIGRATE} down

generate:
	# generate css from tailwind
	@npm run build
	# generate go files from templ templates
	@templ generate

build: generate
	if [ -f "${BINARY}" ]; then \
		rm ${BINARY}; \
		echo "Deleted ${BINARY}"; \
	fi	
	@echo "Building binary..."
	go build -o tmp/${BINARY} cmd/aggregator/*.go 

run: build
	./tmp/${BINARY}

stop:
	@echo "stopping server.."
	@-pkill -SIGTERM -f "./${BINARY}"
	@echo "server stopped..."