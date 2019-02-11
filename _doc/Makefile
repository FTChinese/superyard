build_dir := build
output := backyard_api
input := frontmatter.md admin.md staff.md next-api.md ftc-user stats.md subscription.md

.PHONY: build linux deploy lastcommit clean
pdf : createdir
	pandoc -s --toc --pdf-engine=xelatex -o $(build_dir)/$(output).pdf $(input)

createdir :
	mkdir -p $(build_dir)