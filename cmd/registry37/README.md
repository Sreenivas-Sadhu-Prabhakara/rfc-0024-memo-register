# RFC 37 extension registry checker

RFC 37, Network Meeting Epilogue, etc., by S. D. Crocker (20 March 1970), records the March 17 network meeting, implementation-specification commitments, subsequent protocol criticism, and ideas for leaving room for experiments. It also clarifies process terminology and discusses alternative reconnection approaches. The alternatives are tentative ideas, not a finished protocol specification.

Read and metadata verified 2026-11-28:
https://www.rfc-editor.org/rfc/rfc37.txt
https://www.rfc-editor.org/info/rfc37
https://www.rfc-editor.org/rfc/rfc37.json

Category/status UNKNOWN; source/stream Legacy. No updates, updated-by, obsoletes or obsoleted-by relationships are listed. References to earlier RFCs are context, not formal update relationships.

Modern relevance: reserve extension space deliberately, examine complexity before committing implementations, and separate identifiers from transport-resource allocation. The document raises extension and experimentation concerns without providing a modern negotiated extension mechanism.

Three possible applications inspired by it:
1. Banking API application-error-code catalog linting.
2. Cloud event-type registry collision checks.
3. Distributed-systems experimental opcode allocation exercises.

Selected exercise: an offline Go registry checker for a toy 8-bit code space. A caller supplies stable_max and experimental_min. Stable codes occupy 0..stable_max; experimental codes occupy experimental_min..255; the middle gap is reserved. The example chooses 0..31 and 224..255. These exact boundaries and lifecycle rules are educational choices, not RFC 37 requirements or a production allocation recommendation.

The report flags every affected row for duplicate names or duplicate numerical codes, assignments in reserved space, and mismatches between lifecycle and range. Budgets count distinct occupied code values in each numerical range, including policy violations. Duplicate rows do not consume extra points. Free reserved points remain unavailable under the supplied policy.

Tradeoff: a small explicit registry makes namespace collisions visible, but does not implement negotiation, compatibility testing, migrations, distributed allocation, or public registration. A clean registry proves only these supplied records satisfy the chosen checks. The tool allocates no code externally and sends no network messages.

Run from repository root with Go 1.22+:
```sh
go run ./cmd/registry37 < cmd/registry37/registry.json
go run ./cmd/registry37 -html < cmd/registry37/registry.json > docs/rfc37/index.html
go test -race ./...
go vet ./...
go build ./...
```

Input: one JSON object with a policy and up to 512 entries (empty registry allowed). Both policy boundaries are required, including explicit zero when desired. Require 0 <= stable_max < experimental_min <= 255. Each entry requires a trimmed, case-sensitive name of 1–100 bytes, an explicit integer code 0..255, and lifecycle stable or experimental. A missing/null code is invalid; explicit code zero is valid. Adjacent boundaries permit a zero-capacity reserved gap. A policy always leaves at least one stable and one experimental code.

The CLI rejects unknown fields, trailing JSON, negative or oversized numbers, invalid labels, and input above 64 KiB. Invalidly structured input returns a nonzero exit. Well-formed registries with policy issues produce a report with clean=false and zero exit so they can be inspected; this is not an automatic deployment gate.

Tests cover every code with both lifecycles (512 combinations), duplicate diagnostics, occupancy accounting, zero and boundary policies, empty registries, invalid input and HTML escaping. CI checks generated-report drift.

Hosted static report: https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc37/
No live Go backend is deployed. Run the CLI for custom input. Source is MIT-licensed under the repository license. Standalone repository creation remains pending: the connector lacks a creation operation, the CLI is signed out, and the browser requires GitHub sign-in; actual source is preserved in the existing public repository.

Covered RFC 37; next RFC 38. RFC 26 was previously flagged Not Issued.
