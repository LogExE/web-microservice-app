FROM golang:1.24 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/rest ./cmd/rest
RUN go build -v -o /usr/local/bin/filter ./cmd/filter
RUN go build -v -o /usr/local/bin/notifications ./cmd/notifications

FROM debian:bookworm

COPY --from=builder \
    /usr/local/bin/rest \
    /usr/local/bin/filter \
    /usr/local/bin/notifications \
    /usr/local/bin
