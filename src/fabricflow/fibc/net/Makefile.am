
PY_LIST = fibcnet.py ffpacket.py

.PHONY: go-test py-test

go-test:
	go test -coverprofile=cover.out

py-test:
	pylint ${PY_LIST}
	./suite.py

check-local: go-test py-test
