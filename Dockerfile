# Use the official Go image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download project dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o limeapi main.go

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install necessary certificates and runtime dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/limeapi .

# Copy .env file with a fallback
COPY --from=builder /app/.env .env

# Optional: Set default environment variables
ENV API_PORT=8080
ENV ETH_NODE_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID
ENV DB_CONNECTION_URL=postgresql://user:pass@localhost:5432/dbname
ENV JWT_SECRET=your-secret-key

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./limeapi"]