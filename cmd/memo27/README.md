# RFC 27 memo circulation planner

[Steve Crocker's RFC 27](https://www.rfc-editor.org/info/rfc27/) revised the early Network Working Group's documentation conventions and updated its distribution roster. This educational Go lab uses that change to model a synthetic banking API cutover memo with a changing audience.

## Run from the repository root

Requires Go 1.22 or later.

```sh
go run ./cmd/memo27 < cmd/memo27/revision.json
go run ./cmd/memo27 -html < cmd/memo27/revision.json > report.html
go test -race ./cmd/memo27
```

Input contains a positive revision number, a nonempty roster of unique lowercase site IDs, and a dispatches array (which may be empty). Site IDs have up to 64 ASCII letters or digits with optional hyphen separators. Each dispatch identifies a site and a revision between 1 and the current revision. The input describes a single memo series. It is capped at 64 KiB and rejects unknown fields and trailing JSON.

The planner sorts current sites by ID and picks the highest recorded revision for each. Existing entries for former sites, repeated records, and records received out of order do not create new delivery tasks. The report shows which sites need a first copy, which need an update, and which already have a current dispatch on record. It makes no network requests. Dispatches represent recorded sends, and the plan can change when the ledger or roster is refreshed.

[Hosted sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc27/) · [MIT license](../../LICENSE)
