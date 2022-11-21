FROM --platform=${BUILDPLATFORM} golang:1.19-alpine3.16 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app

# include git
RUN apk add --no-cache git

ARG GOPROXY
# download deps before gobuild
RUN go mod download -x
ARG TARGETOS
ARG TARGETARCH
RUN scripts/container_binary.sh

FROM alpine:3.16

LABEL maintainer="Jeeva Kandasamy <jkandasa@gmail.com>"

ENV APP_HOME="/app" \
    DATA_HOME="/app/_store"

EXPOSE 8080

# install timzone utils
RUN apk --no-cache add tzdata

# create locations
RUN mkdir -p ${APP_HOME} && mkdir -p ${DATA_HOME}

# copy application bin file
COPY --from=builder /app/file-store-server ${APP_HOME}/file-store-server

RUN chmod +x ${APP_HOME}/file-store-server

WORKDIR ${APP_HOME}

ENTRYPOINT [ "/app/file-store-server" ]
CMD [ "-port", "8080" ]
