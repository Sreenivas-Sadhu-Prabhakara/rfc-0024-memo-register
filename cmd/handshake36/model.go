package main

import (
	"fmt"
	"html/template"
	"strings"
)

type Event struct {
	Peer   string `json:"peer"`
	Action string `json:"action"`
}
type Trace struct {
	Name   string  `json:"name"`
	Events []Event `json:"events"`
}
type Scenario struct {
	Traces []Trace `json:"traces"`
}
type State struct {
	Requested bool `json:"requested"`
	Received  bool `json:"peer_request_received"`
	Acked     bool `json:"own_request_acknowledged"`
}

func (s State) ready() bool { return s.Requested && s.Received && s.Acked }

type Step struct {
	Number int `json:"number"`
	Event
	ReadyA    bool   `json:"ready_a"`
	ReadyB    bool   `json:"ready_b"`
	Probe     string `json:"send_probe,omitempty"`
	Violation string `json:"violation,omitempty"`
}
type Result struct {
	Name       string `json:"name"`
	Steps      []Step `json:"steps"`
	A          State  `json:"a"`
	B          State  `json:"b"`
	Violations int    `json:"violations"`
}
type Report struct {
	Results []Result `json:"results"`
}

// A is the fixed send endpoint; B is its receive endpoint.
// Events are observations at the named endpoint, not historical wire commands.
func plan(s Scenario) (Report, error) {
	r := Report{}
	if len(s.Traces) == 0 || len(s.Traces) > 50 {
		return r, fmt.Errorf("expected 1–50 traces")
	}
	seen := map[string]bool{}
	for _, t := range s.Traces {
		if strings.TrimSpace(t.Name) == "" || len(t.Name) > 100 || seen[t.Name] {
			return Report{}, fmt.Errorf("trace names must be nonblank, unique and at most 100 bytes")
		}
		seen[t.Name] = true
		if len(t.Events) == 0 || len(t.Events) > 1000 {
			return Report{}, fmt.Errorf("%s: expected 1–1000 events", t.Name)
		}
		row := Result{Name: t.Name, Steps: []Step{}}
		for i, e := range t.Events {
			if e.Peer != "A" && e.Peer != "B" {
				return Report{}, fmt.Errorf("%s event %d: peer must be A or B", t.Name, i+1)
			}
			local, other := &row.A, &row.B
			if e.Peer == "B" {
				local, other = &row.B, &row.A
			}
			step := Step{Number: i + 1, Event: e}
			switch e.Action {
			case "request":
				local.Requested = true
			case "receive_request":
				if !other.Requested {
					step.Violation = "peer request has not been issued"
				} else {
					local.Received = true
				}
			case "ack":
				if !local.Requested {
					step.Violation = "own request has not been issued"
				} else {
					local.Acked = true
				}
			case "try_send":
				step.Probe = "blocked"
				if e.Peer != "A" {
					step.Violation = "B is the receive endpoint"
				} else if !local.ready() {
					step.Violation = "A is not locally ready"
				} else {
					step.Probe = "allowed"
				}
			default:
				return Report{}, fmt.Errorf("%s event %d: unsupported action", t.Name, i+1)
			}
			if step.Violation != "" {
				row.Violations++
			}
			step.ReadyA = row.A.ready()
			step.ReadyB = row.B.ready()
			row.Steps = append(row.Steps, step)
		}
		r.Results = append(r.Results, row)
	}
	return r, nil
}

var page = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width">
<title>RFC 36 · Handshake readiness</title>
<style>body{font:18px system-ui;max-width:1100px;margin:3rem auto;padding:0 1rem;background:#f4f7fb;color:#172d44}article{margin:2rem 0}table{width:100%;border-collapse:collapse;background:white}th,td{text-align:left;padding:.7rem;border:1px solid #ccd6df}a{color:#075fa3}</style>
<h1>RFC 36 · Handshake readiness</h1>
<p>Educational establishment model: A sends, B receives. Each endpoint becomes locally ready after issuing its own request, receiving the peer request, and observing acknowledgment of its own request.</p>
<p>request = local initiation; receive_request = handling the peer request; ack = handling RFNM for the local request. try_send probes A's establishment readiness only. It does not transmit data or model flow-control permission.</p>
{{range .Results}}<article><h2>{{.Name}}</h2><p>Violations: {{.Violations}}</p><table><thead><tr><th>Step</th><th>Peer / action</th><th>A ready</th><th>B ready</th><th>Send probe</th><th>Violation</th></tr></thead><tbody>
{{range .Steps}}<tr><td>{{.Number}}</td><td>{{.Peer}} / {{.Action}}</td><td>{{.ReadyA}}</td><td>{{.ReadyB}}</td><td>{{.Probe}}</td><td>{{.Violation}}</td></tr>{{end}}</tbody></table></article>{{end}}
<p>Local readiness can differ temporarily. RFNM is not evidence that an application processed data. This deterministic exercise omits NCP encoding, sockets, links, flow control, close, reconnection, loss and timeouts. Repeated observations are idempotent by teaching choice.</p>
<p>Static Go-generated report with no live backend. Historical protocol ideas are used for education; this is not an NCP implementation or a standards-compliance claim.</p>
<p><a href="https://www.rfc-editor.org/rfc/rfc36.txt">Official RFC 36</a> · <a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0024-memo-register/tree/main/cmd/handshake36">Source, tests and instructions</a></p></html>
`))
