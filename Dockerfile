FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download
RUN go get github.com/xuri/excelize/v2

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o overtime-automation ./cmd/overtime/main.go

# Use minimal alpine image for the final container
FROM alpine:latest

# Add SSL certificates for HTTPS connections
RUN apk --no-cache add ca-certificates tzdata

# Copy only the binary from the builder stage
COPY --from=builder /app/overtime-automation /usr/local/bin/overtime-automation

# Set the entrypoint to run the application
ENTRYPOINT ["/usr/local/bin/overtime-automation"]

