# build stage
FROM golang:1.23-alpine AS build-img

WORKDIR /go/src/app
RUN apk update && apk add gcc musl-dev

COPY ./ .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN cd sqlc && sqlc generate

RUN go mod tidy
ENV CGO_ENABLED=1
ENV GOOS=linux
ARG GOARCH
ENV GOARCH=${GOARCH:-amd64}

RUN go test ./...

RUN go build -a -ldflags="-w -s" -o tagyou

# final stage
FROM alpine
EXPOSE 1883 8080
WORKDIR /app
RUN apk update && \
    apk add tzdata sqlite ca-certificates && \
    rm -rf /var/cache/apk && \
    mkdir -p /db
COPY --from=build-img /go/src/app/tagyou /app/
COPY --from=build-img /go/src/app/.env.default.local /app/
CMD ["/app/tagyou"]
