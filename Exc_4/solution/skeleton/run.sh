#!/bin/sh

# todo
# docker build
# docker run db
# docker run orderservice

docker build -t orderservice .

docker rm -f postgres_sbdais
docker rm -f orderservice

./start_postgres.sh

docker run \
    --name orderservice \
    --env-file ./debug.env \
    --network host \
    -d orderservice
