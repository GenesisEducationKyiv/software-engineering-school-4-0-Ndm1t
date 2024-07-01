FROM golang:1.21

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /subscription_service ./cmd/subscription_service

CMD ["/subscription_service"]
