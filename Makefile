BUILD_DIR := build
BINARY := backyard-api

MAC_BIN := $(BUILD_DIR)/mac/$(BINARY)
LINUX_BIN := $(BUILD_DIR)/linux/$(BINARY)

VERSION := `git describe --tags`
BUILD_AT := `date +%FT%T%z`
LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD_AT}"

.PHONY: build run publish linux restart config lastcommit clean test
build :
	go build $(LDFLAGS) -o $(MAC_BIN) -v .

run :
	./$(MAC_BIN)

publish : linux
	rsync -v $(LINUX_BIN) nodeserver:/home/node/go/bin/

linux : 
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(LINUX_BIN) -v .

restart :
	ssh nodeserver supervisorctl restart backyard-api

# Copy env varaible to server
config :
	rsync -v ../.env nodeserver:/home/node/go

lastcommit :
	git log --max-count=1 --pretty=format:%ad_%h --date=format:%Y_%m%d_%H%M

clean :
	go clean -x
	rm -r build/*

test :
	echo $(BUILD)