.PHONY: protobuf
protobuf:
	@grpc_tools_node_protoc --js_out=import_style=commonjs,binary:protobuf --grpc_out=protobuf \
	--plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` hoard.proto