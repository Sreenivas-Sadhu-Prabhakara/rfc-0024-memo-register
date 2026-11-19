# RFC 28 clock accuracy lab

Covered RFC 28. Next RFC 29, Response to RFC 28.

Bill English at SRI asked on 13 January 1970 how accurate a new clock should be for network delay measurements. This short historical discussion request motivates a Go tool that evaluates a latency budget using supplied clock error bounds.

The RFC Editor metadata checked on 19 November 2026 lists Legacy stream, Unknown publication/current status, and empty updates, updated-by, obsoletes and obsoleted-by lists. [Official text and metadata](https://www.rfc-editor.org/info/rfc28/)

## Model

Each timestamp equals true event time plus an error whose absolute value is within the supplied endpoint bound. If the observed timestamp difference is d and the two bounds sum to u, the possible delay is [d-u, d+u], intersected with nonnegative values. If the upper endpoint is negative, the assumptions are inconsistent. These deterministic intervals carry no confidence level.

A delay range wholly at or below the budget is Within budget. A lower endpoint strictly above the budget is Exceeds budget. An interval spanning the budget is Uncertain.

All inputs use integer microseconds on the same reference scale, with no clock changes outside the stated error bounds. Timestamps may be negative and are limited to ±1e12 microseconds, and bounds and budget to 0..1e12, keeping arithmetic in int64 range. Omitted numeric fields take zero. Sample names must be unique nonblank strings of at most 100 bytes. The CLI caps input at 64 KiB and rejects unknown fields and trailing JSON.

This educational model is inspired by RFC 28. Its formula, input schema and verdicts are lab design choices. It measures no real traffic, performs no synchronization and implements no historical wire protocol. It illustrates the cost of uncertainty in banking API telemetry, cloud path comparisons and distributed event ordering. Use a local monotonic clock for elapsed durations within one process. Cross-host timestamp comparisons require known clock relationships.

## Run from the repository root

Requires Go 1.22 or later.

```sh
go run ./cmd/clock28 < cmd/clock28/samples.json
go run ./cmd/clock28 -html < cmd/clock28/samples.json > report.html
go test -race ./cmd/clock28
```

The first two sample requests have the same observed 8000 µs duration. Tight clocks bound it to 6000–10000 µs, while loose clocks yield 0–20000 µs. The smaller uncertainty permits a definite budget assessment. An inconsistent row displays an empty interval in HTML. Its JSON min/max fields retain the raw endpoints and consistent=false.

[Hosted static report](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc28/) · [MIT license](../../LICENSE)
