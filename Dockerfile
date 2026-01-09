FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy source files
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bill-generator .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Create output directory
RUN mkdir -p /app/output

# Copy binary
COPY --from=builder /app/bill-generator .

# Set environment
ENV OUTPUT_DIR=/app/output

# Run
CMD ["./bill-generator"]
