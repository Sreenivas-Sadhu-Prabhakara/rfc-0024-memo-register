# RFC 33 socket label inspector

RFC 33, New Host-Host Protocol (12 February 1970), describes early ARPANET NCP design: intermediary socket naming, simplex connections, flow control, and process lifecycle across heterogeneous hosts. It explicitly lacks enough detail to implement an NCP. Naming abstraction and direction contracts remain useful distributed-systems lessons.

Verified 2026-11-24 against https://www.rfc-editor.org/rfc/rfc33.json and https://www.rfc-editor.org/rfc/rfc33.html:
status/category UNKNOWN; source Legacy; obsoletes RFC 11; updated by RFC 36 and RFC 47; updates none; obsoleted by none listed. This is a historical proposal, not a current Internet Standard.

Three applications inspired by its ideas:
1. Banking adapter endpoint diagnostics in a training tool.
2. Cloud channel send/receive contract inspection.
3. Distributed process naming independently of local process IDs.

Selected exercise: compose symbolic 40-bit labels from a 24-bit user, 8-bit host, and 8-bit AEN. Even AEN receives, odd AEN sends. Check opposite directions and inspect a proposed local AEN relabel, preserving user, host, and AEN parity. A proposal is never applied; invalid direction changes are reported as disallowed. Direction compatibility alone does not establish a connection.

This educational inspector is not an NCP implementation or standards compliance claim. It does not implement wire encoding, transport, control commands, authentication, or flow control. Fixed-width labels make invariants visible but are unsuitable as modern authenticated identities.

Run from repository root with Go 1.22+:
```sh
go run ./cmd/sockets33 < cmd/sockets33/sockets.json
go run ./cmd/sockets33 -html < cmd/sockets33/sockets.json > docs/rfc33/index.html
go test -race ./...
go vet ./...
go build ./...
```

Input: one JSON object with 1–100 uniquely named cases, local and remote socket objects, and optional relabel_aen. Names must be nonblank and at most 100 bytes. Numbers are unsigned decimal integers; omitted numeric fields default to zero. Input above 64 KiB, unknown fields, trailing JSON, and numeric overflow are rejected. Explicit relabel zero differs from omission. Valid-width same-direction cases return opposite_directions=false.

Tests exhaust all 65,536 AEN parity pairs and cover known labels, field bounds, identity preservation, explicit zero, invalid input, and JSON/HTML output. CI checks generated HTML drift.

Hosted output: https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc33/
This is a static Go-generated report; use the CLI for custom input. No live Go backend is deployed.

Source is published under the repository's MIT license. Standalone repository creation remains pending because GitHub requires browser sign-in.

Covered RFC 33; next RFC 34. RFC 26 was previously flagged Not Issued.
