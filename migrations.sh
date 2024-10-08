#!/bin/bash
source .env

if [ -z "${POSTGRES_DB}" ]
then
  POSTGRES_DB=${POSTGRES_USER}
fi

DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
MIGRATIONS_DIR="./migrations"

migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" $@