# RFC 31 symbolic binary fields lab

Covered RFC 31. Next RFC 32, Some Thoughts on SRI's Proposed Real Time Clock.

Daniel Bobrow and William Sutherland's February 1968 paper proposes a symbolic description of binary message bodies. It declares field names, named constants, ways of combining fields, and bit widths. It also discusses alternatives, repetitions, conditions and eventual generation of translation programs. The RFC Editor catalogs this earlier paper as RFC 31, so its date precedes RFCs already covered.

Metadata checked on 22 November 2026 lists Legacy stream, Unknown publication/current status, and no updates, updated-by, obsoletes or obsoleted-by relationships. [Official source](https://www.rfc-editor.org/info/rfc31/)

## Lab scope

This Go teaching tool supports one fixed concatenation of named unsigned fields. Each field emits its least-significant bit first, following the paper's bit-order discussion. The emitted stream fills byte bits 0 through 7, and unused final byte bits are zero. The report shows field offsets, stream bits, packed hexadecimal bytes, and decoded values.

The JSON schema and byte packing convention are lab choices. The historical description permits truncation and sign extension. This lab accepts unsigned values that fit their fields and rejects overflow. It implements a bounded subset of the paper's ideas. Alternatives, repetition, conditional fields and its full description language remain outside the exercise.

Applications include inspecting a banking adapter's binary flags, understanding cloud telemetry layouts, and exploring schema-driven translation. Compact fields save space, while explicit widths, bit order and limits increase the need for precise schemas and validation.

## Run from the repository root

Requires Go 1.22 or later.

```sh
go run ./cmd/fields31 < cmd/fields31/fields.json
go run ./cmd/fields31 -html < cmd/fields31/fields.json > report.html
go test -race ./cmd/fields31
```

Input is a nonempty fields array. Names start with an ASCII letter followed by up to 31 letters, digits or underscores and are unique. Widths are 1..32 bits with at most 512 total bits. Values are unsigned integers below 2^width, with an omitted value meaning zero. The CLI rejects unknown fields, trailing JSON and inputs over 64 KiB.

The sample encodes VERSION=5 in 3 bits, STATUS=17 in 5 bits and REQUEST_ID=257 in 9 bits. Its packed bytes are 8d0101. The internal decoder checks exact byte length and zero padding. Tests include independent golden bytes, maximum values, cross-byte boundaries, truncated and extra payloads, nonzero padding, and invalid layouts.

[Hosted static sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc31/) · [MIT license](../../LICENSE)
