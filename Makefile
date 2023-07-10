config_file_name := api.toml
local_config_file := $(HOME)/config/$(config_file_name)

app_name := superyard
go_version := go1.18.1

current_dir := $(shell pwd)
sys := $(shell uname -s)
hardware := $(shell uname -m)
src_dir := $(current_dir)
out_dir := $(current_dir)/out
build_dir := $(current_dir)/build

default_exec := $(out_dir)/$(sys)/$(hardware)/$(app_name)

linux_x86_exec := $(out_dir)/linux/x86/$(app_name)

linux_arm_exec := $(out_dir)/linux/arm/$(app_name)

.PHONY: build
build :
	go build -o $(default_exec) -tags production -v $(src_dir)

.PHONY: run
run :
	$(default_exec)

.PHONY: builddir
builddir :
	mkdir -p $(build_dir)

.PHONY: devenv
devenv :
	rsync -v $(HOME)/config/env.dev.toml $(build_dir)/$(config_file_name)

.PHONY: version
version :
	git describe --tags > build/version
	date +%FT%T%z > build/build_time

.PHONY: amd64
amd64 :
	@echo "Build production linux version $(version)"
	GOOS=linux GOARCH=amd64 go build -o $(linux_x86_exec) -tags production -v $(src_dir)

.PHONY: arm
arm :
	@echo "Build production arm version $(version)"
	GOOS=linux GOARM=7 GOARCH=arm go build -o $(linux_arm_exec) -tags production -v $(src_dir)

.PHONY: install-go
install-go:
	@echo "Install go version $(go_version)"
	gvm install $(go_version)

.PHONY: config
config :
	mkdir -p ./$(build_dir)
	rsync -v node@tk11:/home/node/config/$(config_file_name) ./$(build_dir)/$(config_file_name)

.PHONY: publish
publish :
	rsync -v $(default_exec) /data/opt/server/API/go/bin/$(app_name)

.PHONY: restart
restart :
	supervisorctl restart superyard

.PHONY: clean
clean :
	go clean -x
	rm build/*

