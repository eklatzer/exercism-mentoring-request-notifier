FROM golang:1.16 as build

ENV APP_HOME /go/src/exercism-mentoring-request-notifier
WORKDIR "$APP_HOME"
COPY . .

RUN go mod download &&\
    go mod verify && \
    go build -x

FROM golang:1.16

ENV APP_HOME /go/src/exercism-mentoring-request-notifier
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=build "$APP_HOME"/exercism-mentoring-request-notifier $APP_HOME

CMD ["./exercism-mentoring-request-notifier"]