APP_NAME = testservice
IMAGE = $(DOCKER_REPO)/$(APP_NAME):$(TAG)
TAG ?= latest

.PHONY: build run docker-build test

build:
	go build -o bin/$(APP_NAME) ./...

run: build
	./bin/$(APP_NAME)

docker-build:
	docker build -t $(IMAGE) .

test:
	go test ./...
