
PYLIST = \
	ofaction.py  \
	ofmatch.py   \
	ofproto.py   \
	offlow.py    \
	ofgroup.py   \
	ofdpa_action \
	ofdpa_match

.PHONY: go-test py-vet

go-test:
	go test -coverprofile=cover.out

py-test:
	pylint --ignore-imports yes $(PYLIST)

check-local: go-test py-test
