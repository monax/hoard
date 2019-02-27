FROM golang:1.11-alpine

ENV PROTO_VERSION 3.3.0

RUN apk add --update --no-cache nodejs npm make curl unzip git python3 bash gcc docker build-base jq

RUN pip3 install docker-compose
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
RUN python3 ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

# install gcloud auth binaries
RUN curl https://sdk.cloud.google.com | bash
RUN ls /root/google-cloud-sdk/bin/
ENV PATH /root/google-cloud-sdk/bin/:$PATH
