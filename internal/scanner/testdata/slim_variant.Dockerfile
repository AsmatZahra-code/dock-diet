# slim_variant.Dockerfile
# Expected: Score=100, Grade=A, Issues=[]
# Verifies: A "-slim" base image (debian:bookworm-slim) is NOT falsely flagged
# as a fat image. Only alpine AND slim images are considered lightweight.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM debian:bookworm-slim AS runner
RUN useradd -m appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
