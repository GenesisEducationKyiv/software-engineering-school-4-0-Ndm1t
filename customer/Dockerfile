FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o customer_service ./cmd/customer

FROM scratch

COPY --from=builder /app/customer_service /app/customer_service

COPY .env .

CMD ["/app/customer_service"]
