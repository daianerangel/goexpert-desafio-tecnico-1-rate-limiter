#!/bin/bash

CONTAINER_NAME="redis-container-test"

# Check if the container exists
if [ "$(docker ps -a --filter "name=$CONTAINER_NAME" --format '{{.Names}}')" == "$CONTAINER_NAME" ]; then
  echo "Container $CONTAINER_NAME exists. Deleting..."
  docker rm -f $CONTAINER_NAME
  echo "Container $CONTAINER_NAME has been deleted."
else
  echo "Container $CONTAINER_NAME does not exist."
fi