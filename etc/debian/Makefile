# -*- coding: utf-8; mode:makefile -*-

MAKESH=./make.sh

# export PROXY=http://192.168.0.1:8080

.PHONY: all
all: modules release


.PHONY: modules
modules:
	$(MAKESH) all


.PHONY: release
release: pkgs frr
	$(MAKESH) release


.PHONY: clean
clean:
	$(MAKESH) clean


.PHONY: pkgs
pkgs:
	$(MAKESH) dl-deb


.PHONY: frr
frr:
	$(MAKESH) dl-frr
