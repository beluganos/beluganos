
PYLIST = \
	fibcevt.py \
	fibccfg.py \
	fibccnv.py \
	fibclog.py

PHONY: go-test py-test

go-test:
	go test -coverprofile=cover.out

py-test:
	pylint --ignore-imports yes $(PYLIST)
	./suite.py

check-local: go-test py-test
