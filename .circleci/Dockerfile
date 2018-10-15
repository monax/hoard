FROM circleci/golang:1.11

RUN curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/protoc-3.3.0-linux-x86_64.zip
RUN unzip protoc-3.3.0-linux-x86_64.zip -d protobuf
RUN sudo cp protobuf/bin/protoc /usr/bin/protoc
RUN rm -rf protobuf protoc-*
RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/goreleaser/goreleaser


# install aws auth binaries
RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
RUN unzip awscli-bundle.zip
RUN sudo ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

# install gcloud auth binaries
RUN curl https://sdk.cloud.google.com | bash
RUN ls /home/circleci/google-cloud-sdk/bin/
ENV PATH /home/circleci/google-cloud-sdk/bin/:$PATH
