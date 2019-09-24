GOIMG=golang:1.13

DOCKERCOMPOSE_PATH=resources/docker-compose/docker-compose.yml
DOCKER_DOCKERFILE=resources/docker/Dockerfile
DOCKER_IMG=nquicenob.com/golang-api-example:local
WORKSPACE=/go/src/nquicenob.com/golang-api-example
DOCKER_BUILD=docker build -f ${DOCKER_DOCKERFILE} -t=${DOCKER_IMG} .

DOCKER_RUN_NAME=golangapi

# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash

.PHONY: build
build:
	docker run \
	--rm  \
	-w $(WORKSPACE) \
	-v ${PWD}:$(WORKSPACE) \
	--env CGO_ENABLED=0 \
	$(GOIMG) go build \
	-a \
	-o build/api \
	cmd/server/main.go

.PHONY: build-image
build-image: build
	$(DOCKER_BUILD)

.PHONY: run
run: 
	docker run \
	-d \
	--rm \
	-w $(WORKSPACE) \
	-v ${PWD}:$(WORKSPACE) \
	--env-file resources/docker-compose/api/api.env \
	--env-file resources/docker-compose/secrets/postgresql.env \
	--name $(DOCKER_RUN_NAME) \
	-p 9000:9000 \
	--network="host" \
	$(GOIMG) go run cmd/server/main.go

.PHONY: build-db-setup
build-db-setup: 
	docker run \
	--rm  \
	-w $(WORKSPACE) \
	-v ${PWD}:$(WORKSPACE) \
	--env CGO_ENABLED=0 \
	$(GOIMG) go build \
	-a \
	-o build/db-setup \
	cmd/db-setup/*.go

.PHONY: dlogs
dlogs: 
	docker logs -f $(DOCKER_RUN_NAME) 

.PHONY: dstop
dstop: 
	docker stop $(DOCKER_RUN_NAME)

clean:
	rm -rf build 

.PHONY: setup
setup: clean build-db-setup build-image 
	./bin/dc-start --build && \
	./bin/dc-wait && \
	./bin/dc-run --entrypoint /db-setup golangapi-dbsetup 

.PHONY: stop
stop: 
	./bin/dc-down
	
.PHONY: test-e2e
test-e2e: 
	./bin/dc-run -d -p 5432:5432 db && \
	./bin/dc-wait && \
	docker run --rm -w $(WORKSPACE) -v ${PWD}:$(WORKSPACE) \
	--env-file resources/docker-compose/api/api.env \
	--env-file resources/docker-compose/secrets/postgresql.env \
	--env-file resources/docker-compose/e2e-test/e2e-test.env \
	--network="host" $(GOIMG) go test -v internal/handlers/*.go || ./bin/dc-down && \
	./bin/dc-down