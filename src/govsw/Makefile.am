SUBDIRS = api pkgs

PACKAGES = govsw/...

.PHONY: go-fmt
go-fmt:
	go fmt ${PACKAGES}

.PHONY: go-vet
go-vet:
	go vet ${PACKAGES}

.PHONY: go-build
go-build:
	go build ${PACKAGES}

.PHONY: gp-install
go-install:
	go install ${PACKAGES}

all-local: go-fmt go-vet go-build

install-exec-local: go-install
