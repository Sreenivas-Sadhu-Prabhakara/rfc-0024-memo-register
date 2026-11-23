# RFC 32 host tracing lab

Covered RFC 32. Next RFC 33, New Host-Host Protocol.

Jerry Cole at UCLA wrote on 5 February 1970 that measurements should include delays inside hosts. His note distinguishes local tracing from absolute one-way delay measurement across hosts, for which clocks require synchronization. It also discusses roughly 1 ms resolution, long-term accuracy, a 24-hour range, and the clock-reading interface.

RFC Editor metadata checked on 23 November 2026 lists Legacy stream, Unknown publication/current status, and no updates, updated-by, obsoletes or obsoleted-by relationships. [Official source](https://www.rfc-editor.org/info/rfc32/)

## Model

A request starts at true time zero, waits in the source host, crosses the network, and is processed by the destination. Both clocks run at the correct rate and resolution is exact. Their constant offsets are configurable.

The source queue span subtracts two source timestamps. The destination service span subtracts two destination timestamps. Each local subtraction cancels that host's offset. Raw network and total delays subtract timestamps from different clocks and include destination_offset - source_offset.

The synthetic truth is supplied as input so the discrepancy can be inspected. The model's constant offsets, exact clocks, event path and JSON schema are educational choices. Drift, clock steps, quantization, return transit and hardware interfaces are outside this small exercise. It performs no synchronization or measurement of real systems.

Applications include banking API latency decomposition, cloud trace debugging, and distributed event timestamp interpretation. Local spans help locate host work. Absolute one-way delay still depends on the relationship between endpoint clocks.

## Run from the repository root

Requires Go 1.22 or later.

```sh
go run ./cmd/trace32 < cmd/trace32/samples.json
go run ./cmd/trace32 -html < cmd/trace32/samples.json > report.html
go test -race ./cmd/trace32
```

All durations are integer microseconds in 0..1e12. Clock offsets are in ±1e12, and omitted numeric fields take zero. Clock pairs need unique nonblank names of at most 100 bytes. Unknown fields, trailing JSON, empty clock arrays and input over 64 KiB are rejected.

The sample has a 5 ms source queue, 8 ms network delay and 10 ms destination service, totaling 23 ms. Local spans remain 5 ms and 10 ms for every pair. When the destination clock is 100 ms behind, raw transit appears as -92 ms and raw total as -77 ms.

[Hosted static sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc32/) · [MIT license](../../LICENSE)
