
PROTOS = ribsapi.proto

.PHONY: proto go-test py-test

go-test:
	go test -coverprofile=cover.out

py-test:
	pylint ${PYLIST}
	./suite.py

proto:
	protoc -I=. --go_out=plugins=grpc:. ${PROTOS}
	protoc -I=. --python_out=. ${PROTOS}

all-local: proto

check-local: go-test py-test
