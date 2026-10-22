#!/bin/bash
docker volume create postgres_sbdais_data
docker run \
    --name postgres_sbdais \
    --env-file ./debug.env \
    -p 5555:5555 \
    --network host \
    -v postgres_sbdais_data:/var/lib/postgresql/18/docker \
    -d postgres:18
