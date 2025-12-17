# Library Usage Guide

This document targets developers embedding the Go module directly (without the CLI). It describes how to load YAML models, run structural validation, execute rules, and customize behaviour.

## 1. Install the module

```
go get github.com/NovokshanovE/archlint
```

Import the packages you need:

```go
import (
    "os"

    "github.com/NovokshanovE/archlint/pkg/archlint"
    "github.com/NovokshanovE/archlint/pkg/engine"
    "github.com/NovokshanovE/archlint/pkg/types"
)
```

## 2. Load and validate a model

```go
f, err := os.Open("path/to/arch.yaml")
if err != nil {
    return err
}
defer f.Close()

model, err := archlint.LoadModelFromYAML(f)
if err != nil {
    return err
}

structuralFindings := archlint.ValidateModel(model)
```

`LoadModelFromYAML` returns a strongly typed `*model.Architecture`. `ValidateModel` performs schema-level checks (version, required fields, duplicates, unknown references). These findings should always be processed first; any `SeverityError` here typically means downstream rules cannot run reliably.

## 3. Run rule engine

```go
opts := engine.Options{} // defaults: all built-in rules enabled
ruleFindings := archlint.RunAll(model, opts)

allFindings := append(structuralFindings, ruleFindings...)
```

Available built-in rules and IDs:

| Rule ID | Purpose |
|---------|---------|
| `ARCH-ACYCLIC` | Detect dependency cycles among containers. |
| `ARCH-CRUD` | Enforce CRUD/database contracts (only tagged services hit DBs, repo-only restrictions). |
| `ARCH-ACL` | Allow external integrations only via ACL-tagged containers. |
| `ARCH-BOUNDARIES` | Report on cohesion vs coupling per boundary (configurable thresholds). |
| `ARCH-EXTERNAL-PROTOCOL` | Require whitelisted protocol prefixes for external calls. |
| `ARCH-DB-ISOLATION` | Ensure databases remain passive and warn on unused databases. |

All rules emit `types.Finding` structures with deterministic ordering.

## 4. Configuring rules programmatically

Populate `engine.Options` to customize execution:

```go
opts := engine.Options{
    EnabledRules: []string{"ARCH-ACYCLIC", "ARCH-ACL"},
    RuleConfig: map[string]map[string]any{
        "ARCH-BOUNDARIES": {
            "minInternalToCrossRatio": 2.0,
            "maxCrossRelations":       10,
        },
        "ARCH-EXTERNAL-PROTOCOL": {
            "allowedPrefixes": []string{"https://gateway.", "amqp://"},
        },
    },
}
```

- `EnabledRules` acts as an allowlist. Leave `nil` to run every registered rule.
- `RuleConfig` forwards arbitrary JSON-like objects to the rule’s decoder (see `pkg/checks/*` for supported fields).

## 5. Loading rule configs from YAML

Instead of building `engine.Options` manually, load them from a YAML file using `pkg/config`:

```go
import "github.com/NovokshanovE/archlint/pkg/config"

opts, err := config.LoadOptionsFromFile("configs/rules.yaml")
if err != nil {
    return err
}
findings := archlint.RunAll(model, opts)
```

Config file schema:

```yaml
rules:
  - id: ARCH-BOUNDARIES
    enabled: true        # omit or set false to skip
    config:
      minInternalToCrossRatio: 1.5
      maxCrossRelations: 5
  - id: ARCH-CRUD
    enabled: false       # disabled
```

## 6. Processing findings

A finding has the structure:

```go
type Finding struct {
    RuleID   string
    Severity types.Severity // "error", "warn", "info"
    Message  string
    Path     string         // e.g. boundaries[0].relations[2]
    Meta     map[string]any // optional, rule-specific context
}
```

Typical integration pattern:

1. Sort or filter findings by `Severity` to decide whether to fail CI or send notifications.
2. Surface `Path` to help users jump to the offending YAML location.
3. If you need text/JSON formatting out of the box, reuse `pkg/report` (`report.WriteText` / `WriteJSON`).

## 7. Extending with custom rules

Rules implement the simple interface in `pkg/checks/checks.go`:

```go
type Rule interface {
    ID() string
    Run(*model.Architecture, map[string]any) []types.Finding
}
```

To add your own rule:

1. Create a file under `pkg/checks` implementing the interface.
2. Register it in `pkg/checks/registry.go` (add to the `rules` slice).
3. Optionally expose configuration options via a struct + `decodeConfig` helper.
4. Add tests and fixtures under `pkg/checks` / `testdata`.

Once registered, the rule becomes available both programmatically and through the CLI/config loader.

## 8. Example: embedding everything

```go
func Evaluate(path string, rulesConfig string) ([]types.Finding, error) {
    fh, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer fh.Close()

    model, err := archlint.LoadModelFromYAML(fh)
    if err != nil {
        return nil, err
    }

    findings := archlint.ValidateModel(model)

    var opts engine.Options
    if rulesConfig != "" {
        opts, err = config.LoadOptionsFromFile(rulesConfig)
        if err != nil {
            return nil, err
        }
    }

    findings = append(findings, archlint.RunAll(model, opts)...)
    return findings, nil
}
```

This function mirrors the intended service use case: ingest YAML, optionally load a rule profile, and return all findings for the caller.
