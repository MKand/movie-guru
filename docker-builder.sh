#!/bin/bash

docker login -u ${DOCKER_NAME} -p ${DOCKER_PASS}
# Get the tag name from the command line argument or use "latest" as default
TAG_NAME=ghack-sre

# Build the Docker image
docker build -t ${DOCKER_NAME}/movie-guru-server:$TAG_NAME ./chat_server_go
docker push ${DOCKER_NAME}/movie-guru-server:$TAG_NAME
echo "Image ${DOCKER_NAME}/movie-guru-server:$TAG_NAME built and pushed successfully!"

# # Build the Docker image
# docker build -t ${DOCKER_NAME}/movie-guru-otelcol:$TAG_NAME ./metrics
# docker push ${DOCKER_NAME}/movie-guru-otelcol:$TAG_NAME
# echo "Image ${DOCKER_NAME}/movie-guru-otelcol:$TAG_NAME built and pushed successfully!"

# # Build the Docker image
# docker build -t ${DOCKER_NAME}/movie-guru-mockuser:$TAG_NAME ./js/mock-user
# docker push ${DOCKER_NAME}/movie-guru-mockuser:$TAG_NAME
# echo "Image ${DOCKER_NAME}/movie-guru-mockuser:$TAG_NAME built and pushed successfully!"

# # Build the Docker image
# docker build -t ${DOCKER_NAME}/movie-guru-frontend:$TAG_NAME ./frontend
# docker push ${DOCKER_NAME}/movie-guru-frontend:$TAG_NAME
# echo "Image ${DOCKER_NAME}/movie-guru-frontend:$TAG_NAME built and pushed successfully!"


# # Build the Docker image
# docker build -t ${DOCKER_NAME}/movie-guru-db:$TAG_NAME ./pgvector
# docker push ${DOCKER_NAME}/movie-guru-db:$TAG_NAME
# echo "Image ${DOCKER_NAME}/movie-guru-db:$TAG_NAME built and pushed successfully!"
