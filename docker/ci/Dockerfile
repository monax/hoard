FROM circleci/golang:1.8

# This Dockerfile is to generate the docker build image for CI services
# See the update_docker_image Make target

RUN curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/protoc-3.3.0-linux-x86_64.zip
RUN unzip protoc-3.3.0-linux-x86_64.zip -d protobuf
RUN sudo cp protobuf/bin/protoc /usr/bin/protoc
RUN rm -rf protobuf protoc-*
RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/goreleaser/goreleaser
