FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o informing_service ./cmd/informing_service

FROM scratch

COPY --from=builder /app/informing_service /app/informing_service

COPY .env .

CMD ["/app/informing_service"]
