CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src/github.com/whosonfirst/go-whosonfirst-sqlite; then rm -rf src/github.com/whosonfirst/go-whosonfirst-sqlite; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-sqlite
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-sqlite/database
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-sqlite/tables
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-sqlite/utils
	cp -r database/* src/github.com/whosonfirst/go-whosonfirst-sqlite/database/
	cp -r tables/* src/github.com/whosonfirst/go-whosonfirst-sqlite/tables/
	cp -r utils/* src/github.com/whosonfirst/go-whosonfirst-sqlite/utils/
	cp -r *.go src/github.com/whosonfirst/go-whosonfirst-sqlite/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/mattn/go-sqlite3"
	@GOPATH=$(GOPATH) go install "github.com/mattn/go-sqlite3"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-index"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-log"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-names"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt database/*.go
	go fmt tables/*.go
	go fmt utils/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-sqlite-index cmd/wof-sqlite-index.go
