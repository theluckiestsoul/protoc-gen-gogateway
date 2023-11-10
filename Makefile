clean:
	sudo rm -rf pkg internal
.PHONY: clean

install:
	cd cmd/protoc-gen-gogateway && go install && cd ../..
.PHONY: install

generate:
	buf generate
.PHONY: generate