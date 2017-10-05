
PYLIST = \
	fibcevt.py \
	fibccfg.py \
	fibccnv.py \
	fibclog.py

py-test:
	pylint --ignore-imports yes $(PYLIST)
	./suite.py

check-local: py-test
