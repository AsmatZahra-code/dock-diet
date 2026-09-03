# perfect.Dockerfile
# Expected: Score=100, Grade=A, Issues=[], NeedsFix=false
# Verifies: A well-optimised Dockerfile passes all rules cleanly.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:3.20 AS runner
RUN adduser -D appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
