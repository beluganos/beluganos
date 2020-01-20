
PROTOS = fibcapi.proto fibcapis.proto
PYLIST = fibcapi.py

.PHONY: proto go-test py-test

go-test:
	go test -coverprofile=cover.out

py-test:
	pylint ${PYLIST}
	./suite.py

proto:  $(PROTOS)
	protoc -I=. --go_out=plugins=grpc:. $<
if ENABLE_PYTHON_GRPC
	python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. $<
endif

all-local: proto

check-local: go-test py-test
