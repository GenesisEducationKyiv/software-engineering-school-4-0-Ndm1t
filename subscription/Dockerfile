FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o subscription_service ./cmd/subscription_service

FROM scratch

COPY --from=builder /app/subscription_service /app/subscription_service

COPY .env .

CMD ["/app/subscription_service"]
