build_dir := build
artifact := backyard-api
doc_file := backyard_api_documentation
inputfiles := doc/frontmatter.md build/doc.md

.PHONY: build linux deploy lastcommit pdf createdir clean
build :
	go build -o $(build_dir)/$(artifact) -v .

deploy : linux
	rsync -v $(build_dir)/$(artifact) nodeserver:/home/node/go/bin/

linux : 
	GOOS=linux GOARCH=amd64 go build -o $(build_dir)/$(artifact) -v .
	
lastcommit :
	git log --max-count=1 --pretty=format:%ad_%h --date=format:%Y_%m%d_%H%M

pdf : createdir
	pandoc -s --toc --pdf-engine=xelatex -o $(build_dir)/$(doc_file).pdf $(inputfiles)

createdir :
	mkdir -p $(build_dir)

clean :
	go clean -x
	rm build/*