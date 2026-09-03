# latest_tag.Dockerfile
# Expected: Score=90, Grade=A, Issues=[Tag], Penalty=-10
# Verifies: :latest tag is detected on the runner stage.
# Note: alpine:latest contains "alpine" so BaseImage rule does NOT fire.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:latest AS runner
RUN adduser -D appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
