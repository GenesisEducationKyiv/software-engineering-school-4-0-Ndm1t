FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o rate_service ./cmd/rate_service

FROM scratch

COPY --from=builder /app/rate_service /app/rate_service

COPY .env .

CMD ["/app/rate_service"]
