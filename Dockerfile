FROM golang:1.13.0-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build

# Let's cache modules retrieval - those don't change so often
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code necessary to build the application
# You may want to change this to copy only what you actually need.
COPY . .

# Build the application
RUN go build -installsuffix cgo -o _output/rest-api cmd/server/main.go

RUN go build -installsuffix cgo -o _output/db-setup cmd/db-setup/main.go


# Create the runtime image
FROM alpine:latest

# add ca-certificates to call https apis
RUN apk --no-cache add ca-certificates

# Create a group and user
RUN addgroup -S akgroup && adduser -S akuser -G akgroup

# Tell docker that all future commands should run as the akuser user
USER akuser

COPY --from=builder /build/_output/rest-api /usr/local/bin/rest-api
COPY --from=builder /build/_output/db-setup /usr/local/bin/db-setup

ENTRYPOINT [ "/bin/sh", "-c", "db-setup && rest-api" ]
