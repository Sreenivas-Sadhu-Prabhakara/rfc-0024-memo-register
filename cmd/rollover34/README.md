# RFC 34 counter rollover lab

RFC 34, Some Brief Preliminary Notes on the Augmentation Research Center Clock, by W. K. English (26 February 1970), describes the ARC XDS 940 clock: calendar fields in two memory words refreshed each second and a 24-bit relative accumulator ticking every 100 microseconds or millisecond. It discusses expected setting accuracy, oscillator stability, standby power and memory-access delays. These are historical hardware observations, not modern timing guarantees.

Verified 2026-11-25 using https://www.rfc-editor.org/rfc/rfc34.txt and https://www.rfc-editor.org/rfc/rfc34.json:
status/category UNKNOWN; source Legacy; no updates, updated-by, obsoletes or obsoleted-by relationships listed.

Modern relevance: distinguish wall-clock representation, counter resolution, clock accuracy, read uncertainty, and rollover. Proposed applications: banking API latency telemetry training; cloud monitoring counter rollover checks; bounded embedded-device trace reconstruction.

Selected exercise: a pure Go inspector for a 24-bit modulo counter. It assumes two readings from one uninterrupted, fixed-rate counter and a caller-supplied inclusive maximum number of counter advances. It reports every possible count through a count and extrema; possibilities are spaced by 2^24 ticks. No measured time, network input, absolute-time decoding, reset detection, drift model or hardware emulation is included.

The formula is d=(end+2^24-start) mod 2^24; possible advances are d+k*2^24 for nonnegative k within the bound. Zero possibilities means inconsistent input assumptions. One means a unique tick count. Several means ambiguous. A unique count is not exact physical elapsed time: phase, quantization and hardware errors remain outside the model. In particular, zero advances does not prove zero elapsed time.

Wrap periods derived from the historical widths and tick rates are 1,677,721,600 microseconds (27m57.7216s) and 16,777,216,000 microseconds (4h39m37.216s). Faster ticks improve resolution but wrap ten times sooner. Extra sampling or an independently justified bound is needed to disambiguate missing wraps. Never infer the bound from the readings themselves.

Run at repository root with Go 1.22+:
```sh
go run ./cmd/rollover34 < cmd/rollover34/samples.json
go run ./cmd/rollover34 -html < cmd/rollover34/samples.json > docs/rfc34/index.html
go test -race ./...
go vet ./...
go build ./...
```

Input is one JSON object containing 1–100 samples. Each sample has a unique nonblank name (at most 100 bytes), start and end (0..16777215), tick_us (100 or 1000), and max_advance_ticks (inclusive, 0..10^12). Omitted unsigned fields default to zero; omitted tick_us is invalid. The artificial 10^12 bound limits scope and keeps duration arithmetic safe. Unknown fields, negative numbers, numeric overflow, input above 64 KiB and trailing JSON are rejected.

Examples include a unique 7-tick crossing, ambiguity after extending the bound by a full wrap, a contradictory bound, and equal readings hiding a wrap. Tests cover boundaries, independent candidate enumeration, bounds, overflow prevention and the CLI.

This is an educational mathematical model inspired by RFC 34. It does not claim historical compliance or recommend the obsolete clock design for production.

Hosted static Go-generated sample: https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc34/
No live backend is hosted. Run the CLI for custom input. Source and tests use the repository's MIT license. Standalone repository creation is pending because no connector operation is available and the browser creation surface failed to initialize in this run; actual source is checked into this existing public repository.

Covered RFC 34; next RFC 35. RFC 26 was previously flagged Not Issued.
