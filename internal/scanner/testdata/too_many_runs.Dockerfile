# too_many_runs.Dockerfile
# Expected: Score=85, Grade=B, Issues=[Layers], Penalty=-15
# Verifies: More than 2 RUN instructions triggers the layer consolidation warning.
# There are 4 RUN instructions: mod download, vet, build, adduser.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go vet ./...
RUN go build -o myapp .

FROM alpine:3.20 AS runner
RUN adduser -D appuser
USER appuser
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
