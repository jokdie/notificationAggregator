FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o provider \
    ./cmd

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/provider .

EXPOSE 8080

ENTRYPOINT ["./provider"]
