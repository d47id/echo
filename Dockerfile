FROM golang:1.11.5-alpine AS builder
ARG VERSION
ARG BUILD_TIME
ARG COMMIT

RUN apk update && apk add build-base git

COPY . /go/src/app
WORKDIR /go/src/app
RUN GO111MODULE=on \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o app \
	-ldflags "-s -w -X lifecycle.Version=${VERSION} -X lifecycle.buildTime=${BUILD_TIME} -X lifecycle.Commit=${COMMIT}"

FROM alpine:3.9
COPY --from=builder /go/src/app/app /
ENTRYPOINT ["/app"]
