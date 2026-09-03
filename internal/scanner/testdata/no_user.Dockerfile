# no_user.Dockerfile
# Expected: Score=80, Grade=B, Issues=[Security], Penalty=-20
# Verifies: A Dockerfile with no USER instruction is flagged as running as root.
# Note: alpine, multi-stage, no apt-get, 1 RUN — all other rules pass.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:3.20 AS runner
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
