# Update the project's version of the hoard protocol with the one at your local go path
.PHONY: update_hoard_proto
update_hoard_proto:
	@scripts/update-hoard-proto.sh

.PHONY: build_protobuf
build_protobuf:
	@grpc_tools_node_protoc --js_out=import_style=commonjs,binary:protobuf --grpc_out=protobuf \
	--plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` hoard.proto