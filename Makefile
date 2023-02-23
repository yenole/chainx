GO ?= go
POD ?=podman

.PHONY: binary
binary: dist
	$(GO) version
	$(GO) build -trimpath --ldflags "-s -w" -o ./dist/chainx ./cmd/chainx/

dist:
	mkdir $@

install:
	$(POD) run -itd --name chainx -p 8085:8080 --env-file $(shell pwd)/dist/.env -v $(shell pwd)/dist/:/data -w /data/ ubuntu:22.04 ./chainx


dev:
	CX_ENV=.env air
