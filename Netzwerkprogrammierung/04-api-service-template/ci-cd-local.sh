#!/bin/bash

set -e

# Define environment variables
export POSTGRES_DB=userdb
export POSTGRES_USER=user
export POSTGRES_PASSWORD=password
export DATABASE_URL=postgres://user:password@localhost:5432/userdb?sslmode=disable

# Define Docker image names
DOCKER_IMAGE_NAME=my-go-app
DOCKER_REGISTRY=your-docker-registry
DOCKER_USERNAME=your-docker-username
DOCKER_PASSWORD=your-docker-password

# Clean up any previous containers
docker compose down -v

# Build stage
echo "Starting build stage..."
docker compose build

# Test stage
echo "Starting test stage..."
docker compose up -d db
sleep 10 # Wait for the database to initialize

docker run --rm --network host -e DATABASE_URL=$DATABASE_URL -v $(pwd):/app -w /app golang:1.22-alpine sh -c "
  apk add --no-cache git bash &&
  go mod download &&
  go test -v -coverprofile=coverage.out ./... &&
  go tool cover -func=coverage.out
"

# Push stage
# echo "Starting push stage..."
# echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin $DOCKER_REGISTRY
# docker tag $DOCKER_IMAGE_NAME $DOCKER_REGISTRY/$DOCKER_IMAGE_NAME:latest
# docker push $DOCKER_REGISTRY/$DOCKER_IMAGE_NAME:latest

# Clean up
docker compose down -v

echo "CI/CD pipeline completed successfully."
