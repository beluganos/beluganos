SUBDIRS = src etc

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

distclean-local:
	-find . -name "Makefile.in" | xargs rm -f
	-rm -f aclocal.m4 config.log config.status
	-rm -fr ./autom4te.cache

fibc-install:
	python ./setup.py install

fibc-egg:
	python ./setup.py bdist_egg --bdist-dir=dist --dist-dir=bin

install-service:
	sudo make -C ./etc/installer install-fibc

uninstall-service:
	sudo make -C ./etc/installer uninstall-fibc

install-stats:
	sudo make -C ./etc/installer install-fibs

uninstall-stats:
	sudo make -C ./etc/installer uninstall-fibs

install-exec-local: fibc-egg

fflow-install:
	install -v -C `pwd`/bin/* ${GOPREFIX}/bin/

release: install fibc-install fflow-install
