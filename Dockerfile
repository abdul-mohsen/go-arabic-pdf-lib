FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY arabictext/ ./arabictext/

# Build from cmd/generator
RUN CGO_ENABLED=0 GOOS=linux go build -o bill-generator ./cmd/generator

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install curl and unzip for font download
RUN apk add --no-cache curl unzip

# Create directories
RUN mkdir -p /app/output /app/fonts

# Download Amiri Arabic font
RUN curl -L -o /tmp/amiri.zip https://github.com/aliftype/amiri/releases/download/1.000/Amiri-1.000.zip && \
    unzip /tmp/amiri.zip -d /tmp/amiri && \
    cp /tmp/amiri/Amiri-1.000/Amiri-Regular.ttf /app/fonts/ && \
    cp /tmp/amiri/Amiri-1.000/Amiri-Bold.ttf /app/fonts/ && \
    rm -rf /tmp/amiri.zip /tmp/amiri

# Copy binary
COPY --from=builder /app/bill-generator .

# Copy invoice data JSON files (Arabic and English)
COPY invoice_data.json .
COPY invoice_data_en.json .

# Set environment
ENV OUTPUT_DIR=/app/output
ENV FONT_DIR=/app/fonts
ENV DATA_FILE=/app/invoice_data.json

# Run
CMD ["./bill-generator"]
