FROM golang:1.17 as build

ENV APP_HOME /go/src/exercism-mentoring-request-notifier
WORKDIR "$APP_HOME"
COPY . .

RUN go mod download &&\
    go mod verify && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -x

FROM alpine:3.15

LABEL org.opencontainers.image.source="https://github.com/eklatzer/exercism-mentoring-request-notifier"

ENV APP_HOME /go/src/exercism-mentoring-request-notifier
WORKDIR "$APP_HOME"

COPY --from=build "$APP_HOME"/exercism-mentoring-request-notifier $APP_HOME
RUN mkdir -p "$APP_HOME"/cfg

CMD ["./exercism-mentoring-request-notifier", "-cache=cfg/cache.json", "-config=cfg/config.json"]