# grade_c_boundary.Dockerfile
# Expected: Score=60 (exactly), Grade=C
# Deductions: BaseImage (-20) + Security (-20) = -40  =>  100 - 40 = 60
# Verifies: Grade C boundary (score >= 60) is correctly assigned.
# Note: ubuntu:22.04 is fat (no alpine/slim), no USER instruction present.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM ubuntu:22.04 AS runner
RUN useradd -m appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
