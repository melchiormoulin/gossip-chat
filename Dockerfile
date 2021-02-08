FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/gossip-chat

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build cmd/gossip-chat.go

# Start fresh from a smaller image
FROM alpine:3.13
RUN apk add ca-certificates
WORKDIR /app

COPY --from=build_base /tmp/gossip-chat/gossip-chat /app/gossip-chat
COPY --from=build_base /tmp/gossip-chat/template/ /app/template/
COPY --from=build_base /tmp/gossip-chat/static/  /app/static/



ENTRYPOINT ["./gossip-chat"]
