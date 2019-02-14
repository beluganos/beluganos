SUBDIRS = net container

PACKAGES = fabricflow/util/...

go-fmt:
	go fmt ${PACKAGES}

go-vet:
	go vet ${PACKAGES}

go-build:
	go build ${PACKAGES}

go-install:
	go install ${PACKAGES}

all-local: go-fmt go-vet go-build

install-exec-local: go-install
