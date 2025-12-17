# ArchLint – Architecture Assurance in Go

ArchLint ships a production-ready Go module that validates service architectures described in YAML, with a reusable library plus an optional CLI.

Looking for Russian documentation? See `README.ru.md`, `docs/examples.ru.md`, and `docs/library.ru.md`.

## Features
- Load architecture definitions from YAML (`examples/payments.yaml`).
- Stress-test the engine with a Spotify-scale streaming service example (`examples/music_streaming.yaml`) that touches every rule.
- Structural validation with precise findings (version check, duplicate containers, unknown relations).
- Built-in rules with stable IDs:
  - `ARCH-ACYCLIC` – detect dependency cycles.
  - `ARCH-CRUD` – guard CRUD/database access semantics (repo/relay style rules).
  - `ARCH-ACL` – enforce ACL-only access to external systems.
  - `ARCH-BOUNDARIES` – cohesion/coupling ratios for boundaries, configurable thresholds.
  - `ARCH-EXTERNAL-PROTOCOL` – ensure integrations hit externals only via approved gateways/transports.
  - `ARCH-DB-ISOLATION` – keep databases passive (no outbound calls, warn on unused DBs).
- Rule configuration via YAML (`configs/rules.yaml`) so callers can enable/disable checks or override per-rule settings.
- Deterministic findings API designed for embedding and further automation.
- Thin CLI wrapper (`cmd/archlint`) for CI usage.
- Go `testing` coverage with fixtures under `testdata/` and golden-ish text output via `pkg/report`.

## YAML model
Architecture source of truth is YAML (versioned). Minimum schema:

```yaml
version: 1
meta:
  owner: platform-team
boundaries:
  - name: Payments
    description: Context information (optional)
    tags: [core]
    containers:
      - name: payments-api
        type: service            # service | database | external
        tags: [acl]
      - name: payments-repo
        type: service
        tags: [repo]
      - name: payments-db
        type: database
        technology: postgres
    boundaries: []               # nested contexts allowed
    relations:
      - from: payments-api
        to: payments-repo
        kind: sync               # sync | async | db
        description: orchestrates domain logic
      - from: payments-repo
        to: payments-db
        kind: db
externals:
  - name: antifraud
    type: external
    description: third-party system
```

Everything lives under boundaries; externals are optional helpers with `type: external`. Each relation tracks a `Path` so findings can point to `boundaries[0].relations[1]` etc.

## Findings contract

```go
type Severity string // "error" | "warn" | "info"
type Finding struct {
    RuleID   string
    Severity Severity
    Message  string
    Path     string // JSONPointer-like path inside YAML
    Meta     map[string]any
}
```

Findings are sorted by `RuleID` then `Path` for deterministic output.

## Library usage

```go
package main

import (
    "os"

    "github.com/PET-dev-projects/ArchLint/pkg/archlint"
    "github.com/PET-dev-projects/ArchLint/pkg/engine"
)

func main() {
    f, _ := os.Open("examples/payments.yaml")
    model, _ := archlint.LoadModelFromYAML(f)

    findings := append(
        archlint.ValidateModel(model),
        archlint.RunAll(model, engine.Options{}),
    )

    // findings is []types.Finding, ready for your logger/CI/etc.
}
```

`engine.Options` allows enabling a subset of rules and passing JSON-like configuration per rule (e.g. relax boundary ratios or add extra CRUD tags).

## CLI

```
go install ./cmd/archlint

archlint check -f examples/payments.yaml \
  --format text \   # pick text or json
  --fail-on error   # threshold (error|warn|info|none)
```

Exit code is non-zero when the selected `--fail-on` severity (default `error`) is met.

See `docs/examples.md` for additional runbook snippets that exercise each built-in rule against the provided fixtures.

### Rule configuration file

Use `--config configs/rules.yaml` (or your own YAML) to control which rules run and pass custom parameters:

```yaml
rules:
  - id: ARCH-ACYCLIC
    enabled: true
  - id: ARCH-BOUNDARIES
    enabled: true
    config:
      minInternalToCrossRatio: 1.5
      maxCrossRelations: 5
```

Each entry references a rule ID; omit or set `enabled: false` to skip it. Any `config` object is forwarded to the rule’s decoder. If no config file is provided, all built-in rules run with their defaults.

## Tests & fixtures
- `testdata/*.yaml` mirror the original PlantUML-based scenarios: cycles, CRUD breaches, ACL violations, weak boundaries.
- `examples/music_streaming.yaml` captures a large streaming platform with multiple boundaries, externals, and data flows so you can validate complex deployments.
- `go test ./...` covers model validation, each rule, and the engine orchestration. Use `GOCACHE=$(pwd)/.cache` if your environment restricts home directories.

## Extending
- Add new rules under `pkg/checks` and register them in `pkg/checks/registry.go`.
- Reuse `pkg/report` for text/JSON output formatting.
- Use `examples/payments.yaml` as a template when migrating from the old PlantUML fixtures.
- For a deeper dive into embedding the library (APIs, rule configuration, extending), see `docs/library.md`.
