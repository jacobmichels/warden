FROM golang:1.19-alpine3.16 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd
COPY internal internal

RUN go build -o warden ./cmd/warden/main.go

FROM alpine:3.16.2

COPY --from=builder app/warden /usr/bin/warden

ENTRYPOINT [ "/usr/bin/warden" ]
