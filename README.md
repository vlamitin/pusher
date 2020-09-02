# Pusher
Connector to pushover

## Prerequisites
### To build and run:
 - docker
 - docker-compose
 
### To develop or run tests:
 - go 1.14
 - golang-migrate cli
 - psql
 - docker
 - docker-compose
 - receive user and app token at pushover.net

## Build
- `cp env.example .env`
- fill .env with your values (specifically PUSHOVER_APP_TOKEN and PUSHOVER_USER_TOKEN)
 - `make build_pusher`
 - `make build_migrator`

## Run
 - `make start`
 - `make stop`

## Test
 - `make test` - run unit tests
 - `make test_e2e` - run e2e tests (need running service)
