FROM golang:1.15-alpine3.12
ENV PROTO_VERSION 3.11.4
ENV GORELEASER_VERSION "v0.149.0"

RUN apk add --update --no-cache \
  nodejs \
  yarn \
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
RUN yarn install -g yarn-cli-login

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

# install gcloud auth binaries
RUN curl https://sdk.cloud.google.com | bash
RUN ls /root/google-cloud-sdk/bin/
ENV PATH /root/google-cloud-sdk/bin/:$PATH

