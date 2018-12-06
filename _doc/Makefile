build_dir := build
output := backyard_api
input := frontmatter.md admin.md staff.md ftcapi.md subscription.md stats.md ftcuser.md

.PHONY: build linux deploy lastcommit clean
pdf : createdir
	pandoc -s --toc --pdf-engine=xelatex -o $(build_dir)/$(output).pdf $(input)

createdir :
	mkdir -p $(build_dir)