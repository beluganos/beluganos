
PYLIST = \
	ofc.py \
	ofdpa2_builtin.py \
	default.py \
	generic.py \
	ovs.py \
	onsl.py \
	ofdpa2.py


py-test:
	pylint --disable=duplicate-code $(PYLIST)
	./suite.py

check-local: py-test
