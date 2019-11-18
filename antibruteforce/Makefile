protogen:
	protoc -I/usr/local/include -I.   -I$(GOPATH)/src   -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis   --go_out=plugins=grpc:. api/api.proto
ipstore:
	go test -v -count=1 antibruteforce/internal/store/ipstore
fmt:
	gofmt -w -s -d .
vet:
	go vet ./...
lint:
	golint ./...
imports:
	goimports -w .
tidy:
	go mod tidy
verify:
	go mod verify

pre-commit: fmt vet lint tidy verify imports
