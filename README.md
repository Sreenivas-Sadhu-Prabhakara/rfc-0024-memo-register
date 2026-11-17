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
