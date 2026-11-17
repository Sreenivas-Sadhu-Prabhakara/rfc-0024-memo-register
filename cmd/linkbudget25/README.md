# RFC 25 link budget lab

[Steve Crocker's October 1969 note](https://www.rfc-editor.org/info/rfc25/) proposed keeping ordinary link numbers at 63 or below so hosts could shrink their tables and leave room for instrumentation. It also mentioned 32 as a possible smaller budget. This educational model compares three direct-index table sizes against the same synthetic ID trace. This is a capacity exercise inspired by ARPANET.

All three policies reserve **255** for diagnostics as an explicit modeling choice. IDs from zero through the ceiling are data, in-range IDs above the ceiling and below 255 are unassigned, and integers outside 0–255 are invalid. Every event is counted, including repeats. Table slots represent relative storage for direct indexing with identical entries. The tool performs no packet parsing or allocation and makes no network requests.

## Run from the repository root

Requires Go 1.22 or later.

```sh
go run ./cmd/linkbudget25 < cmd/linkbudget25/trace.json
go test -race ./cmd/linkbudget25
```

Optional HTML rendering:

```sh
go run ./cmd/linkbudget25 -html < cmd/linkbudget25/trace.json > report.html
```

The input is a single JSON array of integers. Diagnostics and invalid observations are reported as counts. Malformed input is an error. The sample includes values around both smaller ceilings, the diagnostics ID, repeated IDs, and out-of-range values.

License: [MIT](../../LICENSE).
