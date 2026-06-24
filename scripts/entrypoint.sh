#!/bin/sh

# Exit immediately if any command fails
set -e

# If we pass "migrate" as an argument, only do DB creation and migrations
if [ "$1" = "migrate" ]; then
    echo "Checking/Creating database '$POSTGRES_DB'..."
    
    # Try to connect to the target DB. If it fails, attempt to create it using the 'postgres' maintenance DB
    if ! PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1" > /dev/null 2>&1; then
        echo "Database '$POSTGRES_DB' does not exist. Creating via 'postgres' maintenance DB..."
        PGPASSWORD=$POSTGRES_PASSWORD createdb -h "$POSTGRES_HOST" -U "$POSTGRES_USER" "$POSTGRES_DB" || echo "Database might already exist or creation deferred."
    else
        echo "Database '$POSTGRES_DB' already exists and is reachable."
    fi

    echo "Running migrations..."
    ./migrate init || true
    ./migrate up
    echo "Migrations successfully applied!"
    exit 0
fi

# Default behavior for the API container
echo "Starting API..."
exec ./main