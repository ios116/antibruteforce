protogen:
	protoc -I/usr/local/include -I.   -I$(GOPATH)/src   -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis   --go_out=plugins=grpc:. api/api.proto
ipstore:
	go test -v -count=1 antibruteforce/internal/store/ipstore
mockgen:
	mockgen -destination=./internal/usecase/ipusecase/mock_repo.go -package=ipusecase -source=./internal/domain/entities/ip_entities.go IPStoreManager
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
.PHONY: build
build:
	go build -o build/abf
