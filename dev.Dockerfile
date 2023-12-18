# syntax=docker/dockerfile:1.0.0-experimental
FROM golang:1.21-alpine as golang-builder

RUN apk update && apk add --no-cache git openssh-client
RUN mkdir -p -m 0600 /root/.ssh && touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com > /root/.ssh/known_hosts
RUN git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
RUN go env -w GOPRIVATE=github.com/hexcraft-biz/*

WORKDIR /go/src/github.com/hexcraft-biz/drawbridge
COPY . .
RUN --mount=type=ssh go mod tidy
RUN --mount=type=ssh go install ./

FROM alpine
COPY --from=golang-builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=golang-builder /go/bin/drawbridge /var/www/app/
WORKDIR /var/www/app
EXPOSE 9525
ENTRYPOINT /var/www/app/drawbridge
