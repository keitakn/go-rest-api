FROM golang:1.12.5-alpine3.9 as build

WORKDIR /go/app

COPY . .

RUN set -x && \
  apk update && \
  apk add --no-cache git && \
  go build -o go-rest-api

FROM alpine:3.9

WORKDIR /app

COPY --from=build /go/app/go-rest-api .

RUN set -x && \
  addgroup go && \
  adduser -D -G go go && \
  chown -R go:go /app/go-rest-api

CMD ["./go-rest-api"]
