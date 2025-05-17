FROM golang:1.24-alpine AS builder

# Install git (for go modules)
RUN apk add --no-cache git

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o farmako .

# Final image
FROM alpine:latest

# Set workdir in the final image
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/farmako .

# Expose app port
EXPOSE 8080

# Run the app
CMD ["./farmako"]
