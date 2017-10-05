PROTOS = ribsapi.proto
PY_LIST = ribsapi.py

.PHONY: proto go-test py-test

go-test:
	go test -coverprofile=cover.out

py-test:
	pylint ${PY_LIST}

proto:
	protoc -I=. --go_out=plugins=grpc:. ${PROTOS}
	python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ${PROTOS}

all-local: proto

check-local: go-test py-test
