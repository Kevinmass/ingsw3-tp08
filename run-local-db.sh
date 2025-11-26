#!/bin/bash

# Script to start a local PostgreSQL database using Docker for local development.
# This runs PostgreSQL with the same configuration as docker-compose but exposed on localhost:5432.
# Use this as a "mock" database for local testing without affecting the main docker-compose setup.

CONTAINER_NAME="local-postgres-ingw3"
IMAGE="postgres:15-alpine"

echo "Starting local PostgreSQL database..."
echo "If this container already exists, it will be removed first."

# Stop and remove if exists
docker stop $CONTAINER_NAME > /dev/null 2>&1
docker rm $CONTAINER_NAME > /dev/null 2>&1

# Run Postgres with config matching docker-compose
docker run -d \
    --name $CONTAINER_NAME \
    -e POSTGRES_DB=ingsw3_integrated \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=password \
    -p 5432:5432 \
    -v "$(pwd)/db/init.sql:/docker-entrypoint-initdb.d/init.sql" \
    $IMAGE

echo "PostgreSQL is starting on localhost:5432..."
echo "DATABASE_URL for backend: postgresql://postgres:password@localhost:5432/ingsw3_integrated?sslmode=disable"

# Optionally wait for health
echo "Waiting for database to be ready..."
sleep 5
until docker exec $CONTAINER_NAME pg_isready -U postgres > /dev/null 2>&1; do
    echo "Waiting..."
    sleep 2
done

echo "Database is ready!"
