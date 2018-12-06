build_dir := build
BINARY := backyard-api

VERSION := `git describe --tags`
BUILD := `date +%FT%T%z`

LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD}"

.PHONY: build linux deploy lastcommit clean
build :
	go build $(LDFLAGS) -o $(build_dir)/$(BINARY) -v .

run :
	./$(build_dir)/${BINARY}

publish : linux
	rsync -v $(build_dir)/linux/$(BINARY) nodeserver:/home/node/go/bin/ && 

linux : 
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(build_dir)/linux/$(BINARY) -v .

restart :
	ssh nodeserver supervisorctl restart backyard-api

# Copy env varaible to server
config :
	rsync -v ../.env nodeserver:/home/node/go

lastcommit :
	git log --max-count=1 --pretty=format:%ad_%h --date=format:%Y_%m%d_%H%M

clean :
	go clean -x
	rm build/*