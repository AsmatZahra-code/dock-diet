# grade_b_boundary.Dockerfile
# Expected: Score=75 (exactly), Grade=B
# Deductions: Tag (-10) + Cache (-15) = -25  =>  100 - 25 = 75
# Verifies: Grade B boundary (score >= 75) is correctly assigned.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:latest AS runner
RUN apt-get update && apt-get install -y curl && useradd -m appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
