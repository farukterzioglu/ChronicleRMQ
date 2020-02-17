FROM golang:1.13.7-alpine3.11 as builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

RUN mkdir /app
ADD . /app/
WORKDIR /app/chronicle-consumer-rmq

RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk add --no-cache bash
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/chronicle-consumer-rmq/main .

EXPOSE 8800
ENTRYPOINT ["./main"]