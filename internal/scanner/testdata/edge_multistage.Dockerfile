# edge_multistage.Dockerfile
# Expected: Score=100, Grade=A, Issues=[]
# Verifies: A stage alias used as a subsequent FROM base is NOT falsely flagged
# as a fat image. This is the regression test for the smart multi-stage parser.
#
# "FROM base AS runner" uses "base" which is a stage alias defined on line 1,
# not a real image name. The analyzer's validStages map must recognise it.
FROM golang:1.26-alpine AS base
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM base AS runner
RUN adduser -D appuser
USER appuser
COPY --from=base /app/myapp .
ENTRYPOINT ["./myapp"]
