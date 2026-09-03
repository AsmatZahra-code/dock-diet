# no_multistage.Dockerfile
# Expected: Score=90, Grade=A, Issues=[MultiStage], Penalty=-10
# Verifies: A single-FROM Dockerfile is flagged for missing multi-stage build.
# Note: alpine, has USER, 1 RUN, no apt-get — all other rules pass.
FROM golang:1.26-alpine
WORKDIR /app
COPY . .
RUN go build -o myapp .
USER nobody
ENTRYPOINT ["./myapp"]
