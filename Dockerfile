ARG GO_VERSION=1.17.8

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

WORKDIR /var/www

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app ./main.go

FROM alpine:3.15

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /var/www/app .

ADD entrypoint.sh entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./app"]
