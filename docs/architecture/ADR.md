# 📑 Architectural Decision Records (ADRs)

This document records the key architectural decisions taken during the design and development of **Dock-Diet**.

---

## Index of Decisions

- [ADR-001: Go Language & Single Binary CLI Architecture](#adr-001-go-language--single-binary-cli-architecture)
- [ADR-002: Gamified "Diet Score" & Grade Classification System](#adr-002-gamified-diet-score--grade-classification-system)
- [ADR-003: Non-Destructive Auto-Fix Engine (`Dockerfile.optimized`)](#adr-003-non-destructive-auto-fix-engine-dockerfileoptimized)
- [ADR-004: Daemonless Remote Registry Inspection via API](#adr-004-daemonless-remote-registry-inspection-via-api)
- [ADR-005: Tag-Driven Automated CI/CD Release & 1-Line Installer](#adr-005-tag-driven-automated-cicd-release--1-line-installer)
- [ADR-006: Verification Gate with Concurrency Race Detection](#adr-006-verification-gate-with-concurrency-race-detection)

---

### ADR-001: Go Language & Single Binary CLI Architecture

- **Status:** Accepted
- **Context:** CLI users and CI pipelines need a scanner that executes in milliseconds without requiring heavy runtimes (Node.js, Python) or complex system libraries.
- **Decision:** Build `dock-diet` using Go and the `spf13/cobra` framework.
- **Consequences:** 
  - Sub-millisecond startup and execution speeds.
  - Zero runtime dependencies required on target user systems.
  - Easy cross-compilation for Linux, macOS, and Windows.

---

### ADR-002: Gamified "Diet Score" & Grade Classification System

- **Status:** Accepted
- **Context:** Raw lint logs are often ignored by developers. A clear, intuitive metric is needed to gate CI/CD pipelines.
- **Decision:** Implement a deductive scoring engine (0–100) mapping scores to Letter Grades:
  - **Grade A (90–100)**: Excellent optimization
  - **Grade B (75–89)**: Good, minor layer/security fixes suggested
  - **Grade C (60–74)**: Sub-optimal bloat or security risks
  - **Grade D (<60)**: Failing build threshold
- **Consequences:** Pipelines can enforce a `fail_under` threshold via `.dock-diet.yaml`, providing unambiguous PASS/FAIL exit codes (`0` vs `1`).

---

### ADR-003: Non-Destructive Auto-Fix Engine (`Dockerfile.optimized`)

- **Status:** Accepted
- **Context:** Auto-correcting Dockerfiles can inadvertently break application runtime logic if base images or multi-stage steps are altered automatically.
- **Decision:** 
  1. Only apply safe, non-breaking modifications (injecting non-root `USER` and appending `rm -rf /var/lib/apt/lists/*` to `apt-get` commands).
  2. Write output to a separate `Dockerfile.optimized` file rather than modifying the source file directly.
- **Consequences:** Prevents developer code loss and allows developers to inspect proposed diffs before applying them to production.

---

### ADR-004: Daemonless Remote Registry Inspection via API

- **Status:** Accepted
- **Context:** Scanning remote Docker Hub images usually requires `docker pull`, consuming hundreds of megabytes of bandwidth and requiring a running Docker daemon.
- **Decision:** Leverage `google/go-containerregistry` to inspect image manifests and layer metadata directly over HTTP REST APIs.
- **Consequences:** Remote image scanning executes in ~1-2 seconds with zero local Docker daemon dependency.

---

### ADR-005: Tag-Driven Automated CI/CD Release & 1-Line Installer

- **Status:** Accepted
- **Context:** End-users should not need Go installed to run `dock-diet`. Manual releases are error-prone.
- **Decision:**
  1. Automate releases in GitHub Actions triggered by `git push origin v*.*.*`.
  2. Cross-compile 5 native OS/CPU binaries (`linux-amd64`, `linux-arm64`, `darwin-amd64`, `darwin-arm64`, `windows-amd64.exe`).
  3. Publish pre-built Docker containers to GHCR (`ghcr.io/asmatzahra-code/dock-diet`).
  4. Provide a 1-line installation script (`install.sh`) that installs `dock-diet` to `/usr/local/bin`.
- **Consequences:** Zero-friction installation for all developers across all platforms.

---

### ADR-006: Verification Gate with Concurrency Race Detection

- **Status:** Accepted
- **Context:** Prevent regressions, code rot, and hidden concurrency bugs in future pull requests.
- **Decision:** Enforce a standard `make verify` pre-commit and CI gate: `go mod tidy` ➔ `go vet ./...` ➔ `go build` ➔ `go test -race`.
- **Consequences:** Ensures code quality, dependency cleanliness, and race safety before any build is merged or released.
