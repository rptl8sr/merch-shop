#!/bin/bash
set -e

echo "Checking if Postgres is ready at $PG_HOST:$PG_PORT..."

until pg_isready -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER"; do
  echo "Postgres is unavailable - sleeping"
  sleep 3
done

echo "Postgres is ready - running migrations"

export DATABASE_DSN="postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DATABASE}?sslmode=${PG_SSL_MODE}"
echo "Running migrations with DATABASE_DSN=$DATABASE_DSN"
./goose -dir ./migrations postgres "$DATABASE_DSN" up