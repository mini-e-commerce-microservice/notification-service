# Generates mocks for interfaces
INTERFACES_GO_FILES := $(shell find internal -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

PROTOBUF_FILES := $(shell find api/proto -name "*.proto")
PROTOBUF_GEN_FILES := $(PROTOBUF_FILES:api/proto/%.proto=generated/proto/%.pb.go)
generate_protobuf: $(PROTOBUF_GEN_FILES)
	@echo "Generating protobuf files"

$(PROTOBUF_GEN_FILES): generated/proto/%.pb.go: api/proto/%.proto
	@echo "Generating protobuf files $@ for $<"
	buf generate --path $<

generate: generate_mocks generate_protobuf