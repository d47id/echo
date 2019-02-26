NAME=echo
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
COMMIT_SHA=$(shell git rev-parse HEAD)
SHORT_SHA=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

api/api.pb.go: proto/api.proto
	protoc -I proto/ proto/api.proto --go_out=plugins=grpc:api

build: api/api.pb.go
	docker build \
		--tag $(NAME):latest \
		--tag $(NAME):$(SHORT_SHA) \
		--build-arg VERSION=local-dev \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg BRANCH=$(BRANCH) \
		--build-arg COMMIT=$(COMMIT_SHA) .

run: build
	docker run -d --name echo \
	-p 3000:3000 -p 4000:4000 \
	-v $(shell pwd):/stuff \
	echo:$(SHORT_SHA) \
	--config-file /stuff/echo.yaml

all: build

.PHONY: build