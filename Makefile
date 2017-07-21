BINARY=k8guard-discover

VERSION=`git fetch;git describe --tags > /dev/null 2>&1`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

create-hooks:
	ln -s $(CURDIR)/hooks/pre-commit .git/hooks/pre-commit

deps:
	glide install

glide-update:
	glide cc
	glide update

build-docker:
	docker build -t local/k8guard-discover .

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}


clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	go clean

sclean: clean
	rm glide.lock


.PHONY: build
