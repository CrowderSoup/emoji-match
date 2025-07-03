# Variables
BINARY_NAME=emoji-match
VERSION=$(shell git rev-parse --short HEAD)
DOCKER_IMAGE_NAME=emoji-match

# Targets
.PHONY: all build frontend backend docker clean

all: backend

build: backend

frontend:
	npm run build

backend: frontend
	go build -o $(BINARY_NAME) main.go

docker: build
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist
