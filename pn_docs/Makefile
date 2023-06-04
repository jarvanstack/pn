make_dir:=$(shell pwd)
app_name:=$(shell basename $(make_dir))

docsDir := docs

## install: Install docsify and gitbook-summary
.PHONY: install
install:
	go install github.com/jarvanstack/gitbook-summary@latest

## gen: Gen _sidebar.md file
.PHONY: gen
gen:
	cd $(docsDir) && gitbook-summary && cd $(make_dir)

## up: Docker compose up server
.PHONY: up
up:
	docker-compose  -f docker-compose-nginx.yaml up  -d 

## push: Commit and push to remote repo
.PHONY: push
.IGNORE: push
push: gen
	git add .
	git commit -m "update: Auto commit And push"
	git push origin master

## update: Use update_remote.sh to update remote repo
.PHONY: update
update: push
	./update_remote.sh

## serve: Docsify serve in dev env
.PHONY: serve
serve: gen
	# 使用 nginx 代理 docs 目录
	# docker run --rm -it -p 3000:80 -v $(make_dir)/docs:/usr/share/nginx/html:ro nginx
	docsify serve docs

## sync: Sync local repo to cos bucket
.PHONY: sync
sync: push
	coscli sync  docs/ cos://bmft-blog/ -r

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
