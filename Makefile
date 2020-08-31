include .env

PG_CONN_STRING=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

migrate_db_up:
	migrate -database ${PG_CONN_STRING} -path db/migrations up

migrate_db_down:
	migrate -database ${PG_CONN_STRING} -path db/migrations down

start_dev:
	POSTGRES_USER=${POSTGRES_USER} POSTGRES_PASSWORD=${POSTGRES_PASSWORD} POSTGRES_HOST=${POSTGRES_HOST} POSTGRES_PORT=${POSTGRES_PORT} POSTGRES_DB=${POSTGRES_DB} go run cmd/pusher.go -port ${SERVER_PORT} -pushover_app_token ${PUSHOVER_APP_TOKEN} -pushover_user_token ${PUSHOVER_USER_TOKEN}
