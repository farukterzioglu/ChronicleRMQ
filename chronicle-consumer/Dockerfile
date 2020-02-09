FROM golang:1.13.7-alpine3.11 as builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk add --no-cache bash
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
ENTRYPOINT ["./main"]