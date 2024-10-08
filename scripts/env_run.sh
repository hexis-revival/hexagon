if [ -f ./env_run.sh ]; then
  cd ..
fi

source .env

if [ -z "${POSTGRES_DB}" ]
then
  POSTGRES_DB=${POSTGRES_USER}
fi

go run . --db-host 127.0.0.1 \
         --db-port ${POSTGRES_PORT} \
         --db-username ${POSTGRES_USER} \
         --db-password ${POSTGRES_PASSWORD} \
         --db-database ${POSTGRES_DB} \
         --redis-host ${REDIS_HOST} \
         --redis-port ${REDIS_PORT} \
         $@