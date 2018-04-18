SUBDIRS = src

.PHONY: cover gopkg fibc-install fibc-egg

cover:
	@for path in `find . -name cover.out` ; do \
		echo $$path; \
		go tool cover -func=$$path ; \
	done

clean-local:
	-rm -fr ./bin ./pkg ./dist ./build
	-rm -fr ./src/fibc.egg-info/
	-find . -name "*~" | xargs rm -f
	-find . -name "*.pyc" | xargs rm -f
	-find . -name cover.out  | xargs rm -f
	-find . -name cover.log  | xargs rm -f

fibc-install:
	python ./setup.py install

fibc-egg:
	python ./setup.py bdist_egg --bdist-dir=dist --dist-dir=bin

install-service:
	sudo install -v -C `pwd`/etc/systemd/fibcd.conf /etc/beluganos/
	sudo install -v -C `pwd`/etc/systemd/fibcd.service /etc/systemd/system/
	sudo systemctl daemon-reload

uninstall-service:
	sudo rm -v -f /etc/beluganos/fibcd.conf
	sudo rm -v -f /etc/systemd/system/fibcd.service
	sudo systemctl daemon-reload

install-exec-local: fibc-egg

fflow-install:
	install -v -C `pwd`/bin/* ${GOPREFIX}/bin/

release: install fibc-install fflow-install
