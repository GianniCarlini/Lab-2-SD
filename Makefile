grpc:
	export GO111MODULE=on
	go get github.com/golang/protobuf/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0
	export PATH="$PATH:$(go env GOPATH)/bin"
	protoc --go_out=plugins=grpc:proto helloworld.proto

data1:
	go run dataNode1.go


data2:
	go run dataNode2.go

data3:
	go run dataNode3.go	


namenode:
	go run nameNode.go

cliente:
	go run cliente.go