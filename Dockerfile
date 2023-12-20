# build stage
FROM golang:1.21-alpine as build-img

WORKDIR /go/src/app
RUN apk update && apk add --update gcc musl-dev && rm -rf /var/cache/apk/*

COPY ./ .

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
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk add --no-cache tzdata sqlite
RUN mkdir -p /db
COPY --from=build-img /go/src/app/tagyou /app/
COPY --from=build-img /go/src/app/.env.default.local /app/
ENTRYPOINT ["/app/tagyou"]
