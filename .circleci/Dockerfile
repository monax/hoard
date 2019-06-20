FROM golang:1.12.1-alpine3.9
MAINTAINER Monax <support@monax.io>

ENV PROTO_VERSION 3.3.0
ENV GORELEASER_VERSION "v0.104.1"

RUN apk add --update --no-cache \
 nodejs \
 npm \
 netcat-openbsd \
 git \
 openssh-client \
 openssl \
 make \
 bash \
 gcc \
 g++ \
 jq \
 curl \
 docker \
 libffi-dev \
 openssl-dev \
 python-dev \
 py-pip

RUN pip install docker-compose
RUN npm install -g npm-cli-login

RUN curl -OL https://github.com/google/protobuf/releases/download/v${PROTO_VERSION}/protoc-${PROTO_VERSION}-linux-x86_64.zip
RUN mkdir -p protobuf
RUN unzip protoc-${PROTO_VERSION}-linux-x86_64.zip -d protobuf
RUN cp protobuf/bin/protoc /usr/bin/protoc
RUN rm -rf protobuf protoc-*

RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u golang.org/x/lint/golint
ENV GO111MODULE=on

RUN cd /usr/bin && curl -L https://github.com/goreleaser/goreleaser/releases/download/$GORELEASER_VERSION/goreleaser_Linux_x86_64.tar.gz | tar xz goreleaser

# install aws auth binaries
RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
RUN unzip awscli-bundle.zip
RUN python ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

# install gcloud auth binaries
RUN curl https://sdk.cloud.google.com | bash
RUN ls /root/google-cloud-sdk/bin/
ENV PATH /root/google-cloud-sdk/bin/:$PATH

