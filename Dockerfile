FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /wishlist-api ./cmd/app

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /wishlist-api /usr/local/bin/wishlist-api
COPY migrations ./migrations

EXPOSE 8080

CMD ["wishlist-api"]
