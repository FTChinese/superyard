BUILD_DIR := build
BINARY := superyard

LINUX_BIN := $(BUILD_DIR)/linux/$(BINARY)

VERSION := `git describe --tags`
BUILD_AT := `date +%FT%T%z`
LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD_AT}"

.PHONY: build run publish linux restart config lastcommit clean test
build :
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) -v .

run :
	./$(BUILD_DIR)/$(BINARY)

production :
	./$(BUILD_DIR)/$(BINARY) -production

deploy : linux
	rsync -v $(LINUX_BIN) tk11:/home/node/go/bin/
	ssh tk11 supervisorctl restart superyard

linux : 
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(LINUX_BIN) -v .

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