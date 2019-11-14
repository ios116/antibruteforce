protogen:
	protoc -I/usr/local/include -I.   -I$(GOPATH)/src   -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis   --go_out=plugins=grpc:. api/api.proto

ipstore:
	go test -v -count=1 antibruteforce/internal/store/ipstore
