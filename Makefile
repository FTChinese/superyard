BUILD_DIR := build
config_file := api.toml
BINARY := superyard

DEV_OUT := $(BUILD_DIR)/$(BINARY)
LINUX_OUT := $(BUILD_DIR)/linux/$(BINARY)

LOCAL_CONFIG_FILE := $(HOME)/config/$(config_file)

VERSION := `git describe --tags`
BUILD_AT := `date +%FT%T%z`
COMMIT := `git log --max-count=1 --pretty=format:%aI_%h`

LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD_AT} -X main.commit=${COMMIT}"

BUILD_LINUX := GOOS=linux GOARCH=amd64 go build -o $(LINUX_OUT) $(LDFLAGS) -tags production -v .

.PHONY: dev run linux version config deploy build publish clean
# Development
dev :
	go build $(LDFLAGS) -o $(DEV_OUT) -v .

# Run development build
run :
	./$(DEV_OUT)

# Cross compiling linux on for dev.
linux :
	$(BUILD_LINUX)

version :
	echo $(VERSION)

# From local machine to production server
# Copy env variable to server
config :
	rsync -v $(LOCAL_CONFIG_FILE) tk11:/home/node/config

# Test deploy
deploy : linux
	rsync -v $(LINUX_OUT) tk11:/home/node/go/bin/
	ssh tk11 supervisorctl restart superyard

# For CI/CD
build : version
	gvm install go1.15
	gvm use go1.15
	$(BUILD_LINUX)

publish :
	rsync -v $(LINUX_OUT) tk11:/home/node/go/bin/
	ssh tk11 supervisorctl restart superyard

clean :
	go clean -x
	rm build/*

lastcommit :
	git log --max-count=1 --pretty=format:%ad_%h --date=format:%Y_%m%d_%H%M
