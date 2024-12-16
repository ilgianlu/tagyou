# build stage
FROM golang:1.23-alpine AS build-img

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/app
RUN apk update && apk add gcc musl-dev

COPY ./ .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN cd sqlc && sqlc generate

RUN go mod tidy
ENV CGO_ENABLED=1

RUN go test ./...

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags="-w -s" -o tagyou

# final stage
FROM alpine
EXPOSE 1883 1080 8080
WORKDIR /app
RUN apk update && \
    apk add tzdata sqlite ca-certificates && \
    rm -rf /var/cache/apk
COPY --from=build-img /go/src/app/tagyou /app/
CMD ["/app/tagyou"]
