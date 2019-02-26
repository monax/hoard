FROM golang:1.11-alpine

ENV PROTO_VERSION 3.3.0
ENV DCOMP_VERSION 1.23.2

RUN apk add --update --no-cache nodejs npm make curl unzip git python bash gcc docker build-base

RUN curl -L "https://github.com/docker/compose/releases/download/${DCOMP_VERSION}/docker-compose-linux-x86_64" -o /usr/local/bin/docker-compose
RUN curl -OL https://github.com/google/protobuf/releases/download/v${PROTO_VERSION}/protoc-${PROTO_VERSION}-linux-x86_64.zip
RUN unzip protoc-${PROTO_VERSION}-linux-x86_64.zip -d protobuf
RUN cp protobuf/bin/protoc /usr/bin/protoc
RUN rm -rf protobuf protoc-*

RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/golang/dep/cmd/dep

# install aws auth binaries
RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
RUN unzip awscli-bundle.zip
RUN ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

# install gcloud auth binaries
RUN curl https://sdk.cloud.google.com | bash
RUN ls /root/google-cloud-sdk/bin/
ENV PATH /root/google-cloud-sdk/bin/:$PATH
