
PYLIST = \
	fibcapi.py \
	fibcapp.py \
	fibcptm.py \
	fibcofc.py \
	fibcmod.py \
	fibcpkt.py \
	fibcwap.py \
	monapp.py

py-test:
	pylint --ignore-imports yes $(PYLIST)
	./suite.py

check-local: py-test
