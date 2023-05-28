FROM golang:1.19-alpine as golang

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite sqlite gcc musl-dev

ARG API_KEY
ENV API_KEY=$API_KEY

WORKDIR /app
COPY . /app

RUN addgroup -g 1001 build && adduser -D -G build -u 1001 build
RUN chown -R build:build /app

RUN go mod download
RUN go build -o /app/server

USER 1001

CMD ["/app/server"]
