#!/bin/bash

CONTAINER_NAME="local-file-vault-db"

# Start the Docker container
docker compose up -d

# Function to wait for PostgreSQL to be ready
wait_for_postgres() {
  until docker exec -t "$CONTAINER_NAME" pg_isready -U postgres > /dev/null 2>&1; do
    echo "Waiting for PostgreSQL to be ready..."
    sleep 1
  done
  echo "PostgreSQL is ready!"
}

# Wait for PostgreSQL to be ready
wait_for_postgres

# Seed the database
docker cp "$( dirname "$0" )/db-seed.sql" "$CONTAINER_NAME":/
docker exec -t "$CONTAINER_NAME" psql -U postgres -a -f /db-seed.sql

# Create default upload folders
mkdir uploads/shared
mkdir uploads/admin

# Install dependencies
bun install --frozen-lockfile

# Build and start the project
bun run build
nohup bun run start > /dev/null &
