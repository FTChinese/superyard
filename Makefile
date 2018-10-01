build_dir := build
artifact := backyard-api

.PHONY: build linux deploy lastcommit clean
build :
	go build -o $(build_dir)/$(artifact) -v .

deploy : linux
	rsync -v $(build_dir)/$(artifact) nodeserver:/home/node/go/bin/

linux : 
	GOOS=linux GOARCH=amd64 go build -o $(build_dir)/$(artifact) -v .
	
lastcommit :
	git log --max-count=1 --pretty=format:%ad_%h --date=format:%Y_%m%d_%H%M

clean :
	go clean -x
	rm build/*