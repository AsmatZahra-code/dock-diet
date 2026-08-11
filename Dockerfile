# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o dock-diet .

# Stage 2: Create the lightweight final image
FROM alpine:latest
RUN adduser -D appuser
USER appuser

# Copy the binary from the builder stage
COPY --from=builder /app/dock-diet /usr/local/bin/dock-diet

# Run the CLI tool
ENTRYPOINT ["dock-diet"]