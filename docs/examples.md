# Running the CLI Against Built-In Fixtures

Use these snippets to confirm the binary works end-to-end. Every example assumes you are in the project root and already built `./archlint` (`go build ./cmd/archlint`).

## 1. Happy-path architecture

```
./archlint check -f examples/payments.yaml --format text --fail-on error
```

Expected output:

```
No findings
```

This verifies loading, validation, and all default checks over the example YAML shipped with the repo.

## 2. Detecting dependency cycles (ARCH-ACYCLIC)

```
./archlint check -f testdata/arch_cycle.yaml --format text --fail-on error
```

Expected to include a line similar to:

```
ARCH-ACYCLIC	error	boundaries[0].relations[2]	cycle detected: [repo api repo]
```

## 3. CRUD boundary rules (ARCH-CRUD)

```
./archlint check -f testdata/arch_crud_violation.yaml --format text --fail-on error
```

Sample findings:

```
ARCH-CRUD	error	boundaries[0].relations[0]	container api must declare one of [crud repo relay] to access databases
ARCH-CRUD	error	boundaries[0].relations[2]	container cache is restricted to database relations
```

## 4. ACL enforcement (ARCH-ACL)

```
./archlint check -f testdata/arch_acl_violation.yaml --format text --fail-on error
```

Look for:

```
ARCH-ACL	error	boundaries[0].relations[2]	container api must declare one of [acl] to talk to external audit
```

## 5. Boundary cohesion ratios (ARCH-BOUNDARIES)

```
./archlint check -f testdata/arch_boundary_weak.yaml --format text --fail-on warn
```

Because this rule emits warnings, set `--fail-on warn` to see non-zero exit codes. Expected text:

```
ARCH-BOUNDARIES	warn	boundaries[0]	boundary Reporting cohesion/coupling ratio 0.50 below minimum 1.00	{"internal":1,"cross":2,"ratio":0.5}
```

## 6. External integrations via approved transports (ARCH-EXTERNAL-PROTOCOL)

```
./archlint check -f testdata/arch_external_protocol.yaml --format text --fail-on error
```

Expected:

```
ARCH-EXTERNAL-PROTOCOL	error	boundaries[0].relations[0].protocol	protocol "http://public-gateway/antifraud" for external antifraud is not allowed	{"allowedPrefixes":["https://gateway.","kafka://"]}
ARCH-EXTERNAL-PROTOCOL	error	boundaries[0].relations[1]	relation from payments-repo to external antifraud must define protocol
```

## 7. Database isolation (ARCH-DB-ISOLATION)

```
./archlint check -f testdata/arch_db_isolation.yaml --format text --fail-on warn
```

The first finding blocks databases initiating calls, while the second warns about unused databases.

## 8. Automated regression suite

```
GOCACHE=$(pwd)/.cache go test ./...
```

## 9. Custom rule sets via config

```
./archlint check -f testdata/arch_valid.yaml --config configs/rules.yaml --format text
```

This reads `configs/rules.yaml`, enabling only the listed rule IDs and forwarding any inline configuration to each rule.

## 10. Spotify-scale music streaming reference

```
./archlint check -f examples/music_streaming.yaml --format text --fail-on error
```

Expected output:

```
No findings
```

The YAML spans five boundaries (playback, catalog, personalization, monetization, platform) plus multiple externals. Running this scenario confirms the engine scales to large, realistic deployments and that every default rule works together on one topology.

The `GOCACHE=$(pwd)/.cache go test ./...` snippet above executes the Go unit tests that mirror the PlantUML-era scenarios. Each rule has coverage plus engine/validation scaffolding, so failures here usually mean you introduced a regression.
