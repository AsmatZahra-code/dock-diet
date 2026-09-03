# fat_image.Dockerfile
# Expected: Score=80, Grade=B, Issues=[BaseImage], Penalty=-20
# Verifies: A non-alpine/slim runner image (ubuntu:22.04) is flagged as fat.
# Note: No :latest, no apt-get, 2 RUNs, multi-stage, has USER — all other rules pass.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM ubuntu:22.04 AS runner
RUN useradd -m appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
