build_dir := build
BINARY := backyard-api

VERSION := `git describe --tags`
BUILD := `date +%FT%T%z`

LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD}"

doc_file := backyard_api_documentation
inputfiles := _doc/frontmatter.md _doc/admin.md _doc/staff.md _doc/ftcapi.md _doc/subscription.md _doc/stats.md _doc/ftcuser.md

.PHONY: build linux deploy lastcommit pdf createdir clean
build :
	go build $(LDFLAGS) -o $(build_dir)/$(BINARY) -v .

run :
	./$(build_dir)/${BINARY}

deploy : linux
	rsync -v $(build_dir)/linux/$(BINARY) nodeserver:/home/node/go/bin/

# Copy env varaible to server
config :
	rsync -v ../.env nodeserver:/home/node/go

linux : 
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(build_dir)/linux/$(BINARY) -v .
	
lastcommit :
	git log --max-count=1 --pretty=format:%ad_%h --date=format:%Y_%m%d_%H%M

pdf : createdir
	pandoc -s --toc --pdf-engine=xelatex -o $(build_dir)/$(doc_file).pdf $(inputfiles)

createdir :
	mkdir -p $(build_dir)

clean :
	go clean -x
	rm build/*