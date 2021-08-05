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

# For CI/CD
build : version
	gvm install go1.16
	gvm use go1.16
	$(BUILD_LINUX)

syncconfig :
	rsync -v tk11:/home/node/config/$(config_file) ./$(build_dir)

publish :
	rsync -v ./$(build_dir) /home/node/config
	rsync -v $(LINUX_OUT) /home/node/go/bin/
	supervisorctl restart superyard

clean :
	go clean -x
	rm build/*

