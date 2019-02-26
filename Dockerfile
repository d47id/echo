FROM golang:1.11.5-alpine AS base-builder
ARG VERSION
ARG BUILD_TIME
ARG COMMIT
ARG BRANCH

RUN apk update && apk add build-base git && mkdir -p /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download

FROM base-builder as builder
COPY . /src
WORKDIR /src
RUN GO111MODULE=on \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o app \
	-ldflags "-s -w -X github.com/d47id/lifecycle.Version=${VERSION} -X github.com/d47id/lifecycle.BuildTime=${BUILD_TIME} -X github.com/d47id/lifecycle.Branch=${BRANCH} -X github.com/d47id/lifecycle.Commit=${COMMIT}"

FROM alpine:3.9
COPY --from=builder /src/app /
ENTRYPOINT ["/app"]
