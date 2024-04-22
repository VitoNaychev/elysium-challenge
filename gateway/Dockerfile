# Use an official Golang runtime as a parent image
FROM golang:1.21.4-bookworm as builder

WORKDIR /app

ADD ../ /app

WORKDIR /app/gateway/cmd/

RUN go build main.go

CMD ["./main"]