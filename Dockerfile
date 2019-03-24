FROM golang:1.12.1-alpine AS base-builder
RUN apk update && apk add build-base git && mkdir -p /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download

FROM base-builder as builder
ARG VERSION
ARG BUILD_TIME
ARG COMMIT
ARG BRANCH
COPY . /src
WORKDIR /src
RUN GO111MODULE=on \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o app \
	-ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.Branch=${BRANCH} -X main.Commit=${COMMIT}"

FROM alpine:3.9
RUN GRPC_HEALTH_PROBE_VERSION=v0.2.2 && \
	wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
	chmod +x /bin/grpc_health_probe
COPY --from=builder /src/app /
ENTRYPOINT ["/app"]
