#!/bin/bash

source ./scripts/version.sh

# container registry
REGISTRY='quay.io/jkandasa'
PLATFORMS="linux/arm/v6,linux/arm/v7,linux/arm64,linux/amd64"
IMAGE_TAG=${VERSION}

# build and push to quay.io
docker buildx build --push \
  --progress=plain \
  --build-arg=GOPROXY=${GOPROXY} \
  --platform ${PLATFORMS} \
  --file Dockerfile \
  --tag ${REGISTRY}/file-store-server:${IMAGE_TAG} .
