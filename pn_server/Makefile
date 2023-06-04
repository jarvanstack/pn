make_dir:=$(shell pwd)
app_name:=$(shell basename $(make_dir))

# tag 使用 git commit sort 作为 tag
tag:=$(shell git rev-list --count HEAD)

## tidy: Tidy go mod.	
.PHONY: tidy
tidy:
	go mod tidy

## build: Build app
.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/$(app_name) .

## run: Run app
.PHONY: run
run: build
	./bin/$(app_name) --config ./config-dev.yaml

## up: Docker compose up
.PHONY: up
up:
	docker-compose  -f docker-compose.yaml up -d

## docker.build: Docker build
.PHONY: docker.build
docker.build: build
	docker build -t $(app_name):$(tag) -f ./Dockerfile .
	docker tag $(app_name):$(tag) ccr.ccs.tencentyun.com/dengjiawen/$(app_name):$(tag)

## login: Docker repository login
.PHONY: login
login:
	docker login ccr.ccs.tencentyun.com --username=100020377491 

## push: Docker repository push
.PHONY: push
push: docker.build
	docker tag $(app_name):$(tag) ccr.ccs.tencentyun.com/dengjiawen/$(app_name):$(tag)
	docker push ccr.ccs.tencentyun.com/dengjiawen/$(app_name):$(tag)

## pull: Docker repository pull
.PHONY: pull
pull:
	docker pull ccr.ccs.tencentyun.com/dengjiawen/$(app_name):$(tag)
	docker tag ccr.ccs.tencentyun.com/dengjiawen/$(app_name):$(tag) $(app_name):$(tag)
	docker tag $(app_name):$(tag) $(app_name):latest

## update: Commit and push to remote repo, and use update_remote.sh to update remote repo
.PHONY: update
update:
	make push
	./update_remote.sh

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
