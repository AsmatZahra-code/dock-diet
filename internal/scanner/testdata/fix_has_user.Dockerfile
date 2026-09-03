# fix_has_user.Dockerfile
# Input fixture for AutoFix tests: USER instruction already present.
# AutoFix must NOT append a duplicate USER appuser block.
FROM alpine:3.20
COPY app /app
USER nobody
ENTRYPOINT ["/app"]
