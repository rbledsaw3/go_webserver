#!/usr/bin/env bash

CONTAINER_NAME="chirpy_db"

RUNNING=$(docker ps -q --filter "name=${CONTAINER_NAME}")

if [ -n "$RUNNING" ]; then
    echo "Container ${CONTAINER_NAME} is already running."
    read -p "Restart the whole server (Y/n)? " ANSWER
    ANSWER="${ANSWER:-Y}"

    if [[ "$ANSWER" =~ ^[Yy]$ ]]; then
        echo "Shutting down running containers and removing them..."
        docker compose down --volumes --remove-orphans

        echo "Starting up fresh..."
        docker compose up -d db
        docker compose run migrate
        ./build.sh
    else
        echo "Not restarting. Exiting script..."
        exit 0
    fi

else
    echo "Container '$CONTAINER_NAME' is not running."
    echo "Starting DB, running migrations, then starting the server..."

    docker compose up -d db
    docker compose run migrate
    ./build.sh
fi
