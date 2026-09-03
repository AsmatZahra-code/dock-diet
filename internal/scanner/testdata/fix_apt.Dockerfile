# fix_apt.Dockerfile
# Input fixture for AutoFix tests: apt-get without cache cleanup.
# AutoFix should append "&& rm -rf /var/lib/apt/lists/*" to the RUN line.
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y curl
USER nobody
COPY app /app
ENTRYPOINT ["/app"]
