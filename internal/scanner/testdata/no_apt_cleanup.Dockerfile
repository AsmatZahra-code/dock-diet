# no_apt_cleanup.Dockerfile
# Expected: Score=85, Grade=B, Issues=[Cache], Penalty=-15
# Verifies: apt-get without "rm -rf /var/lib/apt/lists/*" is flagged.
# Note: debian:bookworm-slim contains "slim" so BaseImage rule does NOT fire.
# Combined into one RUN to keep layer count at 2.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM debian:bookworm-slim AS runner
RUN apt-get update && apt-get install -y curl && useradd -m appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
