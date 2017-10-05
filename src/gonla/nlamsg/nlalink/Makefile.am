
PROTOS = nlalink.proto

.PHONY: proto go-test

go-test:
	go test -coverprofile=cover.out

proto:
	protoc -I=. --go_out=plugins=grpc:. ${PROTOS}
	protoc -I=. --python_out=. ${PROTOS}

all-local: proto

check-local: go-test
