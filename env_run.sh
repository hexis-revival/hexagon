#!/bin/sh

# Check if parent directory contains `docker-compose.yml`
HAS_DEPLOYMENT_FOLDER=false

if [ -f ../../docker-compose.yml ]; then
  HAS_DEPLOYMENT_FOLDER=true
fi

# Source .env if it exists
if [ -f .env ]; then
  source .env
fi

# Source deployment folder .env if it exists
if ${HAS_DEPLOYMENT_FOLDER}; then
  if [ -f ../../.env ]; then
    source ../../.env
  fi
fi

# Set default .env values
if [ -z ${DATA_PATH+x} ]; then
  DATA_PATH=".data"

  if ${HAS_DEPLOYMENT_FOLDER}; then
    DATA_PATH="../../.data"
  fi
fi

if [ -z ${POSTGRES_HOST+x} ]; then
    POSTGRES_HOST="127.0.0.1"
fi

if [ -z ${POSTGRES_DB+x} ]; then
    POSTGRES_DB=${POSTGRES_USER}
fi

# Run hexagon
go run . --hnet-host ${HNET_HOST} \
         --hnet-port ${HNET_PORT} \
         --hscore-host ${HSCORE_HOST} \
         --hscore-port ${HSCORE_PORT} \
         --data-path ${DATA_PATH} \
         --db-host ${POSTGRES_HOST} \
         --db-port ${POSTGRES_PORT} \
         --db-username ${POSTGRES_USER} \
         --db-password ${POSTGRES_PASSWORD} \
         --db-database ${POSTGRES_DB}