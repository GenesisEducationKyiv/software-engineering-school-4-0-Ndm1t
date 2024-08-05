FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway

FROM scratch

COPY --from=builder /app/gateway /app/gateway

COPY .env .

CMD ["/app/gateway"]
