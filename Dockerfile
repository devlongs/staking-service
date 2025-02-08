# Build stage.
FROM golang:1.20-alpine AS builder

# Install build tools for go-sqlite3.
RUN apk add --no-cache gcc musl-dev

# Enable CGO.
ENV CGO_ENABLED=1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o stakeway-backend ./cmd/api

# Run stage.
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/stakeway-backend .

EXPOSE 8080

CMD ["./stakeway-backend"]
