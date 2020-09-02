include .env

PG_CONN_STRING=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

build_pusher:
	docker build -t pusher:${VERSION} .

build_migrator:
	docker build -t migrator:${VERSION} db/migrator

start:
	docker-compose up -d

stop:
	docker-compose down

start_dev:
	POSTGRES_USER=${POSTGRES_USER} POSTGRES_PASSWORD=${POSTGRES_PASSWORD} POSTGRES_HOST=${POSTGRES_HOST} POSTGRES_PORT=${POSTGRES_PORT} POSTGRES_DB=${POSTGRES_DB} go run cmd/pusher/pusher.go -port ${SERVER_PORT} -pushover_app_token ${PUSHOVER_APP_TOKEN} -pushover_user_token ${PUSHOVER_USER_TOKEN}

migrate_db_up:
	migrate -database ${PG_CONN_STRING} -path db/migrations up

migrate_db_down:
	migrate -database ${PG_CONN_STRING} -path db/migrations down

test:
	go test -v ./...

test_e2e:
	SERVER_HOST=${SERVER_HOST} SERVER_PORT=${SERVER_PORT} go run cmd/pusher_tester/pusher_e2e_tester.go
