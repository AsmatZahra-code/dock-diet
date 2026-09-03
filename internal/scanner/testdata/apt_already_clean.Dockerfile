# apt_already_clean.Dockerfile
# Expected: Score=100, Grade=A, Issues=[]
# Verifies: apt-get WITH "rm -rf /var/lib/apt/lists/*" does NOT produce a
# false-positive Cache issue. The cleanup is present on the same RUN line.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM debian:bookworm-slim AS runner
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/* && useradd -m appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
