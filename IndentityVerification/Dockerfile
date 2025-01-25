FROM golang:1.20-alpine AS builder
LABEL authors="atoncooper"
WORKDIR /ITES
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main

FROM alpine:latest
WORKDIR /ITES
COPY --from=builder /ITES/main .
ENTRYPOINT ["./main"]