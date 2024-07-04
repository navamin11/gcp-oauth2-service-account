#
# STAGE 1: prepare
#
# Start from golang base image
FROM golang:1.22.3-alpine as build

# Updates the repository and installs git
RUN apk update && apk upgrade && apk --no-cache add git && apk --no-cache add tzdata && rm -rf /var/cache/apk/*

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
COPY *.json .

#
# STAGE 2: build
#
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./app ./cmd

#########################################################
#
# STAGE 3: run
#
# The project has been successfully built and we will use a
# lightweight alpine image to run the server
FROM alpine:latest as Dev

# Adds Package to the image
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    apk add --no-cache tzdata && \
    # apk add --no-cache doas && \
    apk add --no-cache bash && \
    apk add --no-cache sudo && \
    rm -rf /var/cache/apk/*

# Copies the binary file from the BUILD container to /app folder
COPY --from=build --chown=support /go/src/app ./app
COPY --from=build --chown=support /go/src/*.json ./

# Runs the binary once the container starts
CMD ./app