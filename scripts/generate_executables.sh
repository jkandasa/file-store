#!/bin/bash

# this script used to generate binary files
# should be executed from the root locations of the repository

source ./scripts/version.sh

BUILD_DIR=builds
BINARY_DIR=binary
# clean builds directory
rm ${BUILD_DIR}/* -rf

# create directories
mkdir -p ${BUILD_DIR}/${BINARY_DIR}

# download dependencies
go mod tidy

function package {
  local PACKAGE_STAGING_DIR=$1
  local BINARY_FILE=$2
  local FILE_NAME=$3
  local FILE_EXTENSION=$4

  mkdir -p ${PACKAGE_STAGING_DIR}

  # echo "Package dir: ${PACKAGE_STAGING_DIR}"
  cp ${BUILD_DIR}/${BINARY_DIR}/${BINARY_FILE} ${PACKAGE_STAGING_DIR}/${FILE_NAME}${FILE_EXTENSION}

  # copy license
  cp LICENSE ${PACKAGE_STAGING_DIR}/LICENSE.txt

  if [[ ${PACKAGE_STAGING_DIR} =~ "windows" ]]; then
    ARCHIVE_NAME="${PACKAGE_STAGING_DIR}.zip"
    zip -r ${BUILD_DIR}/${ARCHIVE_NAME} ${PACKAGE_STAGING_DIR}
  else
    ARCHIVE_NAME="${PACKAGE_STAGING_DIR}.tar.gz"
    tar -czf ${BUILD_DIR}/${ARCHIVE_NAME} ${PACKAGE_STAGING_DIR}
  fi
  rm ${PACKAGE_STAGING_DIR} -rf
}


# platforms to build
PLATFORMS=("linux/arm" "linux/arm64" "linux/386" "linux/amd64" "linux/ppc64" "linux/ppc64le" "linux/s390x" "darwin/amd64" "windows/386" "windows/amd64")

# compile
for platform in "${PLATFORMS[@]}"
do
  platform_raw=(${platform//\// })
  GOOS=${platform_raw[0]}
  GOARCH=${platform_raw[1]}
  package_name_client="store-${GOOS}-${GOARCH}"
  package_name_server="store-server-${GOOS}-${GOARCH}"

  env GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o ${BUILD_DIR}/${BINARY_DIR}/${package_name_client} -ldflags "-s -w $LD_FLAGS" cmd/client/main.go
  if [ $? -ne 0 ]; then
    echo 'an error has occurred. aborting the build process'
    exit 1
  fi

   env GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o ${BUILD_DIR}/${BINARY_DIR}/${package_name_server} -ldflags "-s -w $LD_FLAGS" cmd/server/main.go
  if [ $? -ne 0 ]; then
    echo 'an error has occurred. aborting the build process'
    exit 1
  fi

  FILE_EXTENSION=""
  if [ $GOOS = "windows" ]; then
    FILE_EXTENSION='.exe'
  fi

  package store-${VERSION}-${GOOS}-${GOARCH} ${package_name_client} store ${FILE_EXTENSION}
  package store-server-${VERSION}-${GOOS}-${GOARCH} ${package_name_server} store-server ${FILE_EXTENSION}

done
