##################################################
## BUILDER IMAGE
##################################################

FROM golang:alpine AS builder

ARG BINARY="api"

## Install external dependencies.
RUN apk add upx gcc musl-dev git make

## Copy source files.
WORKDIR /build
COPY ./api/ ./api/
COPY ./.git/ ./.git/

## Install app dependencies.
ENV GO111MODULE=on
WORKDIR /build/api
RUN go version && make install

## Create production binary at '/build/dist/$BINARY'
RUN make build BARGS="-o ../dist/$BINARY"

## Compress binary with UPX.
RUN upx -9 "../dist/$BINARY"


##################################################
## PRODUCTION IMAGE
##################################################

FROM alpine:3.8 as production

ARG BINARY="api"
ARG BUILD_VERSION="unset"
ENV GO_ENV="production"

## Labels:
LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="stevenxie/rgv-api"
LABEL org.label-schema.description="RGV (Reddit Graph Visualizer) API Server"
LABEL org.label-schema.url="https://rgv.stevenxie.me/"
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/rgv"
LABEL org.label-schema.version="$BUILD_VERSION"

## Install dependencies.
RUN apk add ca-certificates

## Copy production artifacts to /api.
COPY --from=builder /build/dist/${BINARY} /usr/bin/${BINARY}

# TODO: Configure API server healthcheck endpoint.
# COPY ./scripts/healthcheck.sh /usr/bin/healthcheck.sh
# ENV ENDPOINT=http://localhost:3000/healthcheck
# HEALTHCHECK --interval=30s --timeout=30s --start-period=10s --retries=1 \
#   CMD [ "healthcheck.sh" ]

## Expose API port.
EXPOSE 3000

## Set entrypoint.
ENV BINARY=$BINARY
ENTRYPOINT "$BINARY"
