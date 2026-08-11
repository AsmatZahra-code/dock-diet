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
- 📊 **The "Diet Score":** A gamified grading system (A to D) that evaluates your container's health based on standards.
- 🤖 **CI/CD & Pipeline Ready:** Supports structured JSON output and customizable pass/fail score thresholds.

---

## 🛠️ Installation

The easiest way to install Dock-Diet is via Go. Simply run this command in your terminal:

```bash
go install [github.com/AsmatZahra-code/dock-diet@latest](https://github.com/AsmatZahra-code/dock-diet@latest)

(Make sure your ~/go/bin directory is added to your system's PATH).

💻 Local CLI Usage
Dock-Diet provides easy-to-use commands for local development.

1. Scan a local Dockerfile:

Bash
dock-diet scan Dockerfile

2. Auto-fix safe optimization issues:

Bash
dock-diet scan Dockerfile --fix

3. Scan a remote Docker image (Requires no Docker Daemon!):

Bash
dock-diet image alpine:latest
dock-diet image ubuntu:latest

4. Generate JSON output (For scripting or jq):

Bash
./dock-diet scan Dockerfile --output json

⚙️ CI/CD Configuration
You can customize Dock-Diet's behavior in your pipelines by creating a .dock-diet.yaml file in the root of your project:

YAML
# .dock-diet.yaml
fail_under: 80     # Pipeline fails if the Diet Score is below 80

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


Built by Asmat Zahra with ❤️ .
