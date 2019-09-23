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

.PHONY: db-setup
db-setup: 
	docker run \
	--rm \
	-w $(WORKSPACE) \
	-v ${PWD}:$(WORKSPACE) \
	--env-file resources/docker-compose/api/api.env \
	--env-file resources/docker-compose/secrets/postgresql.env \
	--network="host" \
	$(GOIMG) go run cmd/db-setup/*.go

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

.PHONY: logs
logs: 
	docker logs -f $(DOCKER_RUN_NAME) 

.PHONY: stop
stop: 
	docker stop $(DOCKER_RUN_NAME)

.PHONY: get
get:
	docker run --rm  -w $(WORKSPACE) -v ${PWD}:$(WORKSPACE) $(GOIMG) go get -u $(filter-out $@,$(MAKECMDGOALS))

.PHONY: mod
mod:
	docker run --rm  -w $(WORKSPACE) -v ${PWD}:$(WORKSPACE) $(GOIMG) go mod $(filter-out $@,$(MAKECMDGOALS))

clean:
	rm -rf build 

.PHONY: test-e2e
test-e2e: clean build-db-setup build-image 
	./bin/dc-start --build && \
	./bin/dc-wait && \
	./bin/dc-run --entrypoint /db-setup golangapi-dbsetup 
	
.PHONY: jia
jia: 
	docker run -d --rm -w $(WORKSPACE) -v ${PWD}:$(WORKSPACE) --network="host" $(GOIMG) go test


	
# && \
# docker run -d --rm -w $(WORKSPACE) -v ${PWD}:$(WORKSPACE) --network="host" $(GOIMG) go test && \
# ./bin/dc-down


