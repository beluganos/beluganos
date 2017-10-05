
PYLIST = \
	ofc.py \
	ofdpa2_builtin.py \
	default.py \
	generic.py \
	ovs.py \
	ofdpa2.py


py-test:
	pylint --disable=duplicate-code $(PYLIST)

check-local: py-test
