FROM golang:1.17 as build

ENV APP_HOME /go/src/exercism-mentoring-request-notifier
WORKDIR "$APP_HOME"
COPY . .

RUN go mod download &&\
    go mod verify && \
    go build -x

FROM golang:1.17

ENV APP_HOME /go/src/exercism-mentoring-request-notifier
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=build "$APP_HOME"/exercism-mentoring-request-notifier $APP_HOME
RUN mkdir -p "$APP_HOME"/cfg

CMD ["./exercism-mentoring-request-notifier", "-cache=cfg/cache.json", "-config=cfg/config.json"]