default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

set_go_path:
	go env -w GOPATH=$HOME/go

install_gen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

gen:
	oapi-codegen -package client spec/spec.yaml > internal/client/gen.go