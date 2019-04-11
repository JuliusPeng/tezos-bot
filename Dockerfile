# build stage
FROM golang:1.12 AS build-env

ENV GO111MODULE=on

WORKDIR  /go/src/github.com/ecadlabs/tezos-bot

COPY go.mod .
COPY go.sum .

RUN go mod download

ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build

# final stage
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.com/ecadlabs/tezos-bot/tezos-bot /app/tezos-bot
ENTRYPOINT ["/app/tezos-bot"]
