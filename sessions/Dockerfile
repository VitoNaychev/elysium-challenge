# Use an official Golang runtime as a parent image
FROM golang:1.21.4-bookworm as builder

WORKDIR /app

ADD ../ /app

RUN go build -o bin/server /app/sessions/cmd/main.go

FROM debian:bookworm-slim

COPY --from=builder /app/bin/server .

EXPOSE 8080

CMD ["./server"]