
PYLIST = \
	fibcdbm.py \
	fibcmap.py \
	portmap.py \
	dps.py \
	idmap.py


py-test:
	pylint --disable=duplicate-code $(PYLIST)
	./suite.py

check-local: py-test
