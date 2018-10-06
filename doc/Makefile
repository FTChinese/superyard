build_dir := build
outfile := cms_api_documentation
inputfiles := src/metadata.md ../README.md src/login.md src/cms-user.md src/admin.md src/ftc-api.md src/search.md src/data.md

.PHONY: pdf createdir clean
pdf : createdir
	pandoc -s --toc --resource-path=.:img --pdf-engine=xelatex -o $(build_dir)/$(outfile).pdf $(inputfiles)

createdir :
	mkdir -p $(build_dir)

clean :
	rm -r build/*