# Dock-Diet Makefile

APP_NAME = dock-diet

# ── Dependency management ─────────────────────────────────────────────────────

# Ensure go.mod and go.sum are in sync with the source code.
tidy:
	go mod tidy

# ── Code quality ──────────────────────────────────────────────────────────────

# Run the built-in Go static analyser across all packages.
vet:
	go vet ./...

# ── Build ─────────────────────────────────────────────────────────────────────

# Compile the CLI binary into the project root.
build: tidy
	go build -o $(APP_NAME) .

# ── Testing ───────────────────────────────────────────────────────────────────

# Run the full test suite (verbose output).
test:
	go test ./... -v

# Run the full test suite with the Go race condition detector enabled.
# Use this before every commit to catch concurrency bugs early.
test-race:
	go test ./... -v -race

# Run tests, generate a coverage profile, and produce an HTML report.
# Open coverage.html in your browser to see line-by-line coverage.
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report written to coverage.html"

# ── Verification gate (professional pre-commit / CI command) ──────────────────

# verify runs the full professional check sequence in strict order:
#   1. tidy   — keep go.mod / go.sum clean
#   2. vet    — catch common mistakes with the built-in linter
#   3. build  — ensure the code compiles without errors
#   4. test-race — run all tests with the race detector
#
# Any failure stops the pipeline immediately. This mirrors the verification
# gate used by professional Go CLIs (kubectl, helm, goreleaser, etc.).
verify: tidy vet build test-race
	@echo "✅ All checks passed."

# Alias so CI pipelines can call 'make ci' and get the same behaviour as 'make verify'.
ci: verify

# ── Utilities ─────────────────────────────────────────────────────────────────

# Remove the compiled binary and any auto-generated optimized Dockerfiles.
clean:
	rm -f $(APP_NAME)
	rm -f *.optimized

# Build the binary and immediately scan the project's own Dockerfile.
run-scan: build
	./$(APP_NAME) scan Dockerfile