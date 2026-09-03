# fix_no_user.Dockerfile
# Input fixture for AutoFix tests: missing USER instruction.
# AutoFix should append: RUN useradd -m appuser / USER appuser
FROM alpine:3.20
COPY app /app
ENTRYPOINT ["/app"]
