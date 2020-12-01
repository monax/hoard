# We use a multistage build to avoid bloating our deployment image with build dependencies
FROM golang:1.15-alpine3.12 as builder
RUN apk add --update --no-cache make git bash

ARG REPO=$GOPATH/src/github.com/monax/hoard
COPY . $REPO
WORKDIR $REPO

ENV GO111MODULE=on
# Build purely static binaries
RUN make build

# This will be our base container image
FROM alpine:3.12

# We like it when TLS works
RUN apk add --no-cache ca-certificates

ARG REPO=/go/src/github.com/monax/hoard

# Copy binaries built in previous stage
COPY --from=builder $REPO/bin/* /usr/local/bin/

EXPOSE 53431

ENV HOARD_JSON_CONFIG '{"ListenAddress":"tcp://:53431","Storage":{"StorageType":"memory","AddressEncoding":"base64"},"Logging":{"LoggingType":"logfmt","Channels":["info","trace"]}}'

CMD [ "hoard", "--env-config", "-l" ]
