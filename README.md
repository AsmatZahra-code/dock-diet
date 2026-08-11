# 🥗 Dock-Diet

![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)

**Dock-Diet** is a blazingly fast, lightweight CLI tool built in Go. It acts as a fitness coach for your Dockerfiles and container images, helping you reduce size bloat, enforce security best practices, and consolidate layers. 

Whether you are running it locally on your machine or inside a CI/CD pipeline, Dock-Diet ensures your containers stay lean and secure.

---

## 🚀 Key Features

- 🔍 **Advanced Dockerfile Scanning:** Catches bad practices like missing multi-stage builds, running as root, and leftover `apt-get` caches.
- 🛠️ **Auto-Healer (`--fix`):** Safely auto-corrects minor issues and generates a clean `Dockerfile.optimized`.
- ☁️ **Remote Image Scanning:** Analyzes image sizes and layer counts directly from Docker Hub **without** having to `docker pull` them.
- 📊 **The "Diet Score":** A gamified grading system (A to D) that evaluates your container's health based on CNCF standards.
- 🤖 **CI/CD & Pipeline Ready:** Supports structured JSON output and customizable pass/fail score thresholds.

---

## 🛠️ Installation

Clone the repository and build the binary using the provided Makefile:

```bash
git clone [https://github.com/AsmatZahra-code/dock-diet.git](https://github.com/AsmatZahra-code/dock-diet.git)
cd dock-diet
make build
(This will generate a dock-diet executable binary in your root folder).

💻 Local CLI Usage
Dock-Diet provides easy-to-use commands for local development.

1. Scan a local Dockerfile:

Bash
./dock-diet scan Dockerfile
2. Auto-fix safe optimization issues:

Bash
./dock-diet scan Dockerfile --fix
3. Scan a remote Docker image (Requires no Docker Daemon!):

Bash
./dock-diet image alpine:latest
./dock-diet image ubuntu:latest
4. Generate JSON output (For scripting or jq):

Bash
./dock-diet scan Dockerfile --output json
⚙️ CI/CD Configuration
You can customize Dock-Diet's behavior in your pipelines by creating a .dock-diet.yaml file in the root of your project:

YAML
# .dock-diet.yaml
fail_under: 80     # Pipeline fails if the Diet Score is below 80
ignore_rules: []   # Future implementation
🤖 GitHub Action Integration
Dock-Diet is designed to be a native GitHub Action. You can integrate it directly into your repositories to block unoptimized Dockerfiles from being merged.

Add this step to your .github/workflows/ci.yml:

YAML
name: Docker Optimization Check
on: [push, pull_request]

jobs:
  dock-diet-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4
        
      - name: Run Dock-Diet Scanner
        uses: AsmatZahra-code/dock-diet@main
        with:
          dockerfile_path: 'Dockerfile'
🏗️ Project Architecture (Standard Go Layout)
This project strictly follows the Standard Go Project Layout used by CNCF projects:

cmd/ - Contains the CLI entry points (Cobra framework).

internal/scanner/ - The core business logic (analyzers, fixers, remote image parsing).

Makefile - Developer commands for building and testing.

👨‍💻 Developer Guide
If you want to contribute to the code, use these handy Make commands:

make build : Cleans dependencies and builds the dock-diet binary.

make test  : Runs all unit tests locally.

make clean : Removes binaries and temporary .optimized files.

make tidy  : Cleans up go.mod and unused dependencies.

Built by Asmat Zahra with ❤️ .