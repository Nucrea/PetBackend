FROM golang:1.22-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/shortlinks cmd/shortlinks
COPY pkg pkg
COPY internal internal

RUN GOEXPERIMENT=boringcrypto go build -ldflags "-s -w" -o ./app ./cmd/shortlinks
RUN chmod +x ./app

FROM alpine:3.21.2 AS production
WORKDIR /backend

COPY --from=builder /build/app .
COPY deploy/shortlinks-config.yaml ./config.yaml

EXPOSE 8080

CMD ["./app", "-c", "config.yaml"]