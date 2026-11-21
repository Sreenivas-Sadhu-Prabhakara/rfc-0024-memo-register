# RFC 30 memo directory checker

Covered RFC 30. Next RFC 31, Binary Message Forms in Computer.

Steve Crocker's 4 February 1970 revision of the early Network Working Group memo conventions preserves the invitation to share tentative ideas and updates the circulation roster and contact directory. It includes Carnegie-Mellon and Stanford in the recipient list.

RFC Editor metadata checked on 21 November 2026 lists Legacy stream and Unknown status. RFC 30 updates RFCs 10, 16, 24 and 27. Its updated-by, obsoletes and obsoleted-by lists are empty. [Official source](https://www.rfc-editor.org/info/rfc30/)

## Exercise

This modern educational directory checker joins a distribution roster with a separately maintained contact directory by site ID. It derives a useful administrative task from the historical note. Its JSON schema and consistency rules are lab design choices.

A matching contact produces a Ready route. A missing directory record or an exact contact-name mismatch produces an unresolved entry with no ready mailbox. Directory sites outside the roster are reported separately. Names are compared exactly, including case, so a cosmetic name change can require review. Results are sorted by site ID and preserve the source records.

Useful applications include banking API change notices, cloud incident contact reviews and open-source maintainer circulation. The program produces a preview of routing data. Mailbox syntax validation uses Go's net/mail parser and requires a bare address, with no display name. It does not verify mailbox existence, deliverability, or ownership and sends no messages.

## Run from the repository root

Requires Go 1.22 or later.

```sh
go run ./cmd/directory30 < cmd/directory30/directory.json
go run ./cmd/directory30 -html < cmd/directory30/directory.json > report.html
go test -race ./cmd/directory30
```

Input contains a nonempty roster array and a directory array that may be empty. Site IDs have up to 64 lowercase ASCII letters and digits with optional hyphen separators. Names contain up to 100 bytes, have nonblank text, and have no leading or trailing whitespace. Mailboxes are at most 254 bytes. Duplicate site IDs within either input list, unknown fields, trailing JSON and input over 64 KiB fail validation. Valid input exits successfully even when the report identifies unresolved entries.

The fictional sample has one ready route, one changed contact, one missing directory record, and one unused directory site. Example addresses use the reserved .invalid domain.

[Hosted static report](https://sreenivas-sadhu-prabhakara.github.io/rfc-0024-memo-register/rfc30/) · [MIT license](../../LICENSE)
