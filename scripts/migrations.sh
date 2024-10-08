#!/bin/bash
source .env

if [ -z "${POSTGRES_DB}" ]
then
  POSTGRES_DB=${POSTGRES_USER}
fi

DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
MIGRATIONS_DIR="./migrations"

if ! [ -x "$(command -v migrate)" ]; then
  echo 'Error: go-migrate is not installed.' >&2
  echo 'Please follow their installation guide: https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md#installation' >&2
  exit 1
fi

migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" $@