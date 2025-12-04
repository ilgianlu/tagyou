# build build-dev stage
FROM golang:1.25-alpine AS build-img

ARG TARGETOS
ARG TARGETARCH
ARG SQLC_URL=https://downloads.sqlc.dev/sqlc_1.30.0_linux_amd64.tar.gz

RUN apk update && apk add gcc musl-dev curl
RUN curl -L "${SQLC_URL}" | tar -xzvf - -C /usr/bin
RUN ls -la /usr/bin

WORKDIR /go/src/app
COPY ./ .

RUN cd sqlc && \
    sqlc generate && \
    go mod tidy

ENV CGO_ENABLED=1

RUN go test ./...

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags="-w -s" -o tagyou

# final stage
FROM alpine

RUN apk update && \
    apk add tzdata sqlite ca-certificates && \
    rm -rf /var/cache/apk

WORKDIR /app
COPY --from=build-img /go/src/app/tagyou /app/

EXPOSE 1883 1080 8080
CMD ["/app/tagyou"]
