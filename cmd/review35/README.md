# RFC 35 review readiness checker

RFC 35, Network Meeting, by S. D. Crocker (3 March 1970), is an administrative notice. It announces a UCLA meeting on 17 March at 9 a.m. about implementation details following RFC 33, asks for one programmer per site, reserves 18 March for continued discussion, and plans documentation circulation on 20 March if discussion does not supersede it. The notice does not establish that those future events occurred.

Official sources read and metadata verified 2026-11-26:
- https://www.rfc-editor.org/rfc/rfc35.txt
- https://www.rfc-editor.org/info/rfc35/
- https://www.rfc-editor.org/rfc/rfc35.json

Category/status: INFORMATIONAL. Stream/source: Legacy. No updates, updated-by, obsoletes, or obsoleted-by relationships are listed. Its reference to RFC 33 is contextual, not a formal update relationship.

Modern lessons: bring implementing teams into design review, identify a current draft, allow discussion to change it, and circulate documentation deliberately.

Three applications inspired by the notice:
1. Banking API design-review preparation across API, Security, and Operations.
2. Cloud migration review representation checks.
3. Distributed-systems workshop readiness after a draft changes.

Selected exercise: an offline Go checker. For each supplied meeting it reports missing expected sites, sites with multiple delegates, unexpected sites, delegates whose reviewed revision does not match the current draft, and unresolved blocker labels. A meeting is ready to circulate only when all five lists are empty. These are explicit educational policy choices; RFC 35 specifies no machine-readable review process. The tool sends no messages, approves no design, and publishes no document.

Run from repository root with Go 1.22+:
```sh
go run ./cmd/review35 < cmd/review35/reviews.json
go run ./cmd/review35 -html < cmd/review35/reviews.json > docs/rfc35/index.html
go test -race ./...
go vet ./...
go build ./...
```

Input: one JSON object containing 1–50 meetings, each with a unique name, draft_revision, 1–100 unique sites, up to 200 delegates, and up to 100 blocker labels. Labels are exact, case-sensitive, trimmed strings of 1–100 bytes. Each delegate has site, name, and reviewed_revision; an empty or omitted reviewed_revision is reported as unreviewed. An empty delegate list is valid and produces missing sites. Names identify records only, without authentication or verification of attendance. Collection order is preserved for meetings/blockers; site diagnostic lists are sorted. Duplicate delegate rows count as multiple delegates, so input duplication remains visible.

Unknown JSON fields, trailing JSON values, input over 64 KiB, invalid labels, and invalid collection sizes are rejected. Valid records failing readiness return a report with decision hold; operational input errors produce a nonzero exit. This distinction supports inspection rather than an automatic deployment gate.

Tests exercise each hold reason independently, deterministic diagnostics, current-vs-stale revisions, invalid input, JSON/HTML rendering, and HTML escaping. CI checks the generated static report for drift.

Tradeoff: exact labels and a simple one-delegate rule make missing representation visible, but do not model legitimate multi-role attendance, design correctness, quorum policies, waivers, real approvals, or review evidence. Ready means only the supplied data passes these exercise rules. This educational tool is not standards compliance or a production governance recommendation.

Hosted output: https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc35/
This is a static Go-generated report, with no live backend. Run the CLI for custom input. Source is MIT-licensed under the repository license. Standalone repository creation remains pending: the connector lacks a create-repository operation, the CLI is signed out, and the browser requires GitHub sign-in; source is preserved in the existing public repository.

Covered RFC 35; next RFC 36. RFC 26 was previously flagged Not Issued.
