build_dir := build
outfile := ftc_restful_api_documentation
inputfiles := src/metadata.md ../README.md src/article.md src/account.md src/user.md

.PHONY: pdf createdir image clean
pdf : createdir image
	pandoc -s --toc --resource-path=.:img --pdf-engine=xelatex -o $(build_dir)/$(outfile).pdf $(inputfiles)

createdir :
	mkdir -p $(build_dir)

image :
	sh chart.sh

clean :
	rm -r build/*