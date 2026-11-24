# RFC 24 memo register

RFC 24 revised the early Network Working Group's memo-writing and circulation procedures. Its dispatch rule called for one copy per listed recipient from the author's site, leaving further reproduction to recipients.

This Go lab applies that rule to a fictional payments API cutover memo. It is an offline educational audit of recorded dispatches. A ledger entry represents a recorded copy and supplies no receipt confirmation. The lab makes no network requests.

RFC Editor metadata lists RFC 24 under Legacy with status Unknown. It updates RFCs 10 and 16, is updated by RFCs 27 and 30, and has no listed obsoleting relationships. [RFC 24 metadata and text](https://www.rfc-editor.org/info/rfc24)

## Run locally

Requires Go 1.22 or later.

```sh
go run ./cmd/memo24 < examples/distribution.json
go run ./cmd/memo24 -html < examples/distribution.json > report.html
```

The input is one JSON object with a title, a roster of unique recipient IDs, and a copies array with one entry per recorded dispatch. IDs use lowercase letters and digits, with hyphens separating words. An empty copies array records no dispatches. The output lists roster entries in input order and then unexpected recipients in alphabetical order. Counts describe recipients in each category. The example has one matched, one missing, one duplicate, and one unexpected recipient.

The command exits successfully for any valid register, including registers with discrepancies, so its output can feed other reporting tools. Invalid input returns an error.

## Verification

```sh
go test -race ./...
go vet ./...
go build ./...
go run ./cmd/memo24 -html < examples/distribution.json | diff -u docs/index.html -
```

The HTML report under `docs/` is generated from the example. [Open the hosted sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/).

MIT licensed. Next in the series is [RFC 25, No High Link Numbers](https://www.rfc-editor.org/info/rfc25).

## RFC 25 companion lab

[Link budget lab](cmd/linkbudget25) compares three direct-index table sizes against one synthetic trace. [Open the hosted sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc25/).

## RFC 27 companion lab

[Memo circulation planner](cmd/memo27) shows who still needs a copy of the current revision. [Open the hosted sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc27/).

## RFC 28 companion lab

[Clock accuracy lab](cmd/clock28) checks a latency budget using timestamp error bounds. [Open the static sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc28/).

## RFC 29 companion lab

[Clock resolution lab](cmd/resolution29) compares identical synthetic RTTs across three clock tick sizes. [Open the static sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc29/).

## RFC 30 companion lab

[Memo directory checker](cmd/directory30) reconciles roster contacts with directory records before planning a circulation. [Open the static sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc30/).

## RFC 31 companion lab

[Symbolic binary fields lab](cmd/fields31) encodes named unsigned fields, shows offsets and emitted bits, and decodes the packed bytes. [Open the static sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc31/).

## RFC 32 companion lab

[Host tracing lab](cmd/trace32) compares local host spans with timestamp differences across unsynchronized clocks. [Open the static sample](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc32/).

## RFC 33 — Socket label inspector

Covered RFC 33; next RFC 34. [Go source, tests, and run instructions](cmd/sockets33/) · [Hosted static report](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc33/).

Educational naming and direction inspector; no NCP implementation or live backend. RFC 33 has Unknown status, obsoletes RFC 11, and is updated by RFCs 36 and 47. Standalone repository creation remains pending GitHub browser sign-in; source is preserved here under MIT.
