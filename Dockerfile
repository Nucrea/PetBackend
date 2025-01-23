FROM golang:1.22-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/backend cmd/backend
COPY pkg pkg
COPY internal internal

RUN GOEXPERIMENT=boringcrypto go build -ldflags "-s -w" -o ./app ./cmd/backend
RUN chmod +x ./app

FROM alpine:3.21.2 AS production
WORKDIR /backend

COPY --from=builder /build/app .
COPY cmd/backend/config.yaml .
COPY cmd/backend/jwt_signing_key .

EXPOSE 8080

CMD ["./app", "-c", "config.yaml"]