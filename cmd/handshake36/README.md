# RFC 36 handshake readiness lab

RFC 36, Protocol Notes, by S. D. Crocker (16 March 1970), develops early ARPANET host protocol ideas around establishment, receiver-controlled suspension/resumption, and dynamic reconnection. It separates switching from initial establishment, describes NCP bookkeeping tables, and proposes control commands. These historical designs are not production recommendations.

Read and verified 2026-11-27:
https://www.rfc-editor.org/rfc/rfc36.txt
https://www.rfc-editor.org/info/rfc36/
https://www.rfc-editor.org/rfc/rfc36.json

Category/status UNKNOWN; stream/source Legacy. Updates RFC 33; updated by RFC 39 and RFC 44; no obsoletes or obsoleted-by relationships listed.

Three applications inspired by its ideas:
1. Banking integration readiness tests before allowing a request through.
2. Cloud service startup-handshake trace inspection.
3. Distributed-systems concurrency training with reordered observations.

Selected exercise: a fixed, one-way connection establishment model with sender A and receiver B. Each becomes locally ready only after its local request, receipt of its peer's request, and handling of the RFNM acknowledging its own request. Receipt of a peer request can precede local initiation; it is retained in the model. Request delivery and RFNM handling can occur in different orders.

Actions at the named peer:
- request: issue the local request.
- receive_request: handle the peer request; requires the peer to have issued one.
- ack: handle acknowledgment of the local request; requires local initiation.
- try_send: inspect A's establishment readiness. B probes are rejected because B is the receive endpoint.

Input: one JSON object with 1–50 uniquely named traces, each containing 1–1000 events. Names must be nonblank and at most 100 bytes. Peer must be A or B and action must be one of the four above. Unknown fields, trailing JSON, input over 64 KiB and invalid action/peer labels are input errors. Causally impossible observations and premature send probes are reported as violations; they do not invent missing events or change establishment state. Such validly structured traces still produce reports and a zero process exit. Unsupported input produces a nonzero exit.

The two hosted examples each block one premature probe, then allow A's probe after both observations arrive. They also show A becoming ready before B. Repeated establishment observations are idempotent by explicit teaching choice. No payload is transmitted. This tool does not model flow-control permission, wire commands, socket or link allocation, close, reconnection, timeouts, failures, or complete NCP semantics. RFNM is not proof that an application has processed data.

Tradeoff: a deterministic model exposes ordering errors reproducibly, but cannot validate a live network stack. Historical wire compliance and modern production protocol safety are outside its scope.

Run at repository root with Go 1.22+:
```sh
go run ./cmd/handshake36 < cmd/handshake36/traces.json
go run ./cmd/handshake36 -html < cmd/handshake36/traces.json > docs/rfc36/index.html
go test -race ./...
go vet ./...
go build ./...
```

Tests exhaust all 80 causally valid orderings of six distinct establishment observations, verify readiness at every prefix, and cover premature probes, asymmetric readiness, impossible observations, duplicates, malformed input and HTML escaping. CI checks generated HTML drift.

Hosted static report: https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc36/
No live backend is deployed. Run the Go CLI for custom input. Source is MIT-licensed under the repository license. Standalone repository creation remains pending: the connector lacks a creation operation, the CLI is signed out, and the browser requires GitHub sign-in; actual source is preserved in the existing public repository.

Covered RFC 36; next RFC 37. RFC 26 was previously flagged Not Issued.
