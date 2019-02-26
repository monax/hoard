export SRC_PATH := ../protobuf
export DST_PATH := ./protobuf

PROTO_FILES = $(shell find ../ -path ../vendor -prune -o -path ../hoard-js -prune -o -path ../node_modules -prune -o -type f -name '*.proto' -print)

.PHONY: clean
clean:
	rm -f ${DST_PATH}/*pb.js
	rm -f ${DST_PATH}/github.com/gogo/protobuf/gogoproto/*pb.js

# Currently, Hoard-JS only supports dynamic protobuf loading.
# This target is left to exemplify how we may statically compile the proto definintions. 
.PHONY: protobuf
protobuf:
	protoc --proto_path=${SRC_PATH} --js_out=import_style=commonjs,binary:${DST_PATH} -I ${SRC_PATH} \
		--grpc_out=${DST_PATH} --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` ${PROTO_FILES}