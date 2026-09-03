# 🏛️ Dock-Diet Architecture & System Design

This document details the architecture, component flow, and release pipeline of **Dock-Diet**.

> 📖 For key technical choices and design rationale, see [Architectural Decision Records (ADRs)](ADR.md).

---

## 1. High-Level System Architecture

```mermaid
graph TB
    subgraph TriggerLayer ["1. Trigger & Entry Layer"]
        CLI_User["💻 Terminal User<br/><code>dock-diet scan Dockerfile --fix</code>"]
        GHA_User["⚡ GitHub Action Runner<br/><code>action.yaml</code>"]
        Tag_Push["🏷️ Git Tag Push<br/><code>git push origin v1.0.0</code>"]
    end

    subgraph CommandLayer ["2. CLI Command Dispatcher (cmd/)"]
        RootCmd["<code>root.go</code><br/>Cobra Command Root"]
        ScanCmd["<code>scan.go</code><br/>Scan Dockerfile Command"]
        ImageCmd["<code>image.go</code><br/>Remote Image Command"]
    end

    subgraph ScannerCore ["3. Core Engine (internal/scanner/)"]
        ConfigLoader["<code>config.go</code><br/>Load <code>.dock-diet.yaml</code> Thresholds & Rules"]
        Analyzer["<code>analyzer.go</code><br/>Static Analysis & Score Engine"]
        Fixer["<code>fixer.go</code><br/>Auto-Fix Injector"]
        RemoteInspector["<code>image.go</code><br/>Go-ContainerRegistry API Inspector"]
    end

    subgraph OutputLayer ["4. Execution Outputs"]
        Report["📊 Console Output<br/>Diet Score (0-100) & Grade (A/B/C/D)"]
        OptFile["🛠️ Generated File<br/><code>Dockerfile.optimized</code>"]
        ExitCode["🛑 Exit Code<br/>0 = Pass | 1 = Fail Threshold"]
    end

    subgraph ReleasePipeline ["5. CI/CD Release Pipeline (.github/workflows/test-action.yml)"]
        VerifyJob["🧪 Verify Job<br/><code>make verify</code> (tidy, vet, test -race)"]
        BuildBinaries["📦 Cross-Compiler<br/>Linux / macOS / Windows Binaries"]
        GHRelease["🚀 GitHub Release<br/>Upload Assets & Auto-Changelog"]
        GHCRPush["🐳 GHCR Registry<br/>Push <code>ghcr.io/asmatzahra-code/dock-diet</code>"]
    end

    %% Flow Connections
    CLI_User --> RootCmd
    GHA_User --> ScanCmd
    
    RootCmd --> ScanCmd
    RootCmd --> ImageCmd

    ScanCmd --> ConfigLoader
    ScanCmd --> Analyzer
    
    ConfigLoader -.->|Config Overrides| Analyzer
    Analyzer -->|If --fix enabled| Fixer
    Fixer --> OptFile

    ImageCmd --> RemoteInspector

    Analyzer --> Report
    Analyzer --> ExitCode
    RemoteInspector --> Report

    Tag_Push --> VerifyJob
    VerifyJob --> BuildBinaries
    BuildBinaries --> GHRelease
    BuildBinaries --> GHCRPush
```

---

## 2. Dockerfile Scanner & Scoring Workflow

```mermaid
sequenceDiagram
    autonumber
    actor User as User / CI Pipeline
    participant CLI as cmd/scan.go
    participant Config as internal/scanner/config.go
    participant Engine as internal/scanner/analyzer.go
    participant Fixer as internal/scanner/fixer.go
    participant Output as Console / Terminal

    User->>CLI: dock-diet scan Dockerfile --fix
    CLI->>Config: LoadConfig() (Read .dock-diet.yaml)
    Config-->>CLI: Config Options (FailUnder, IgnoreRules)
    
    CLI->>Engine: Analyze(dockerfilePath, config)
    
    Note over Engine: Rule Checks:<br/>1. Base image (slim/alpine)<br/>2. Multi-stage build<br/>3. RUN consolidation<br/>4. apt-get clean & rm -rf<br/>5. Non-root USER instruction
    
    Engine-->>CLI: AnalysisResult (Score, Issues, Grade, NeedsFix)
    
    alt --fix flag is passed AND NeedsFix is true
        CLI->>Fixer: AutoFix(dockerfilePath)
        Fixer->>Fixer: Inject USER & append apt cleanup
        Fixer-->>Output: 🛠️ Created Dockerfile.optimized
    end

    CLI->>Output: 📊 Print Diet Score & Grade (A/B/C/D)
    
    alt Score < Target FailUnder Score
        CLI-->>User: ❌ Exit Code 1 (Pipeline Failed)
    else Score >= Target FailUnder Score
        CLI-->>User: ✅ Exit Code 0 (Pipeline Passed)
    end
```

---

## 3. CI/CD Release Automation Flow

```mermaid
flowchart LR
    subgraph Triggers ["Triggers"]
        PushCode["Git Push to Branch"]
        PushTag["Git Tag Push"]
    end

    subgraph Job1 ["Job 1: Verify"]
        V1["go mod tidy"] --> V2["go vet ./..."]
        V2 --> V3["go build"]
        V3 --> V4["go test -race"]
    end

    subgraph Job2 ["Job 2: Integration Test"]
        I1["Build Binary"] --> I2["Run End-to-End Smoke Tests"]
    end

    subgraph Job3 ["Job 3: Release (On Version Tags)"]
        R1["Extract Version Tag & Lowercase Repo Name"]
        R2["Cross-Compile 5 Binaries<br/>(Linux, Darwin, Windows)"]
        R3["Softprops Release Action<br/>(Upload Assets + Auto-Changelog)"]
        R4["Docker Build & Push to GHCR<br/>(ghcr.io/asmatzahra-code/dock-diet)"]
        
        R1 --> R2 --> R3 --> R4
    end

    PushCode --> Job1
    Job1 --> Job2
    
    PushTag --> Job1
    Job2 -->|Passes| Job3
```
