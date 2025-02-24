FROM golang:1.22-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/notifyer cmd/notifyer
COPY pkg pkg
COPY internal internal

RUN go build -ldflags "-s -w" -o ./app ./cmd/notifyer
RUN chmod +x ./app

FROM alpine:3.21.2 AS production
WORKDIR /backend

COPY --from=builder /build/app .
COPY deploy/notifyer-config.yaml ./config.yaml

EXPOSE 8081

CMD ["./app", "-c", "config.yaml"]