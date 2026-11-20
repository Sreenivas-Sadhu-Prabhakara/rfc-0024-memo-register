# RFC 29 clock resolution lab

Covered RFC 29. Next RFC 30, Documentation Conventions.

Robert Kahn at BBN replied on 19 January 1970 to SRI's question about network clocks. His note expected round trips in the tens to hundreds of milliseconds and described a 16-bit IMP hardware clock ticking every 100 microseconds, with most measurements using a 25.6 ms software clock.

RFC Editor metadata checked on 20 November 2026 lists Legacy stream, Unknown publication/current status, and no updates, updated-by, obsoletes or obsoleted-by entries. [Official source](https://www.rfc-editor.org/info/rfc29/)

## Model and tradeoff

This educational Go simulation measures identical synthetic round trips with tick sizes of 100, 1000, and 25600 microseconds. It computes floor((start+rtt)/tick)*tick - floor(start/tick)*tick. Each elapsed duration uses two readings of the same ideal clock. Both timestamps share a reference scale and a perfect clock rate. Counters are unbounded. Counter wrap, scheduling delays and drift are outside this model. The RFC supplies the clock scales, while flooring and fixture behavior are lab design choices.

A coarser tick loses detail for fast calls. Error can be positive or negative depending on the start phase, with absolute error strictly below a tick. A zero reading can cover a positive duration.

Useful applications include banking API telemetry resolution checks, cloud path comparisons and tracing instrumentation design. The lab highlights the information loss when timestamps are stored or collected at coarse resolution. It supplies synthetic observations for learning.

## Run from the repository root

Requires Go 1.22 or later.

```sh
go run ./cmd/resolution29 < cmd/resolution29/samples.json
go run ./cmd/resolution29 -html < cmd/resolution29/samples.json > report.html
go test -race ./cmd/resolution29
```

Input is one JSON object with a nonempty samples array. Each sample has a unique nonblank name (up to 100 bytes), start_us and rtt_us. Both integers are between zero and 1e12 microseconds. Omitted numeric fields take zero. The CLI rejects unknown fields, trailing JSON and input over 64 KiB. Reports preserve input order.

The default fast API call lasts 800 µs. It is measured as 800 µs, 1000 µs, and 0 µs with the three clock scales. The 15 ms service call also disappears on the 25.6 ms grid at its chosen start phase.

[Hosted static report](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc29/) · [MIT license](../../LICENSE)
