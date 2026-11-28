package main

import (
	"fmt"
	"html/template"
	"strings"
)

type Policy struct {
	StableMax       *uint16 `json:"stable_max"`
	ExperimentalMin *uint16 `json:"experimental_min"`
}
type Entry struct {
	Name      string  `json:"name"`
	Code      *uint16 `json:"code"`
	Lifecycle string  `json:"lifecycle"`
}
type Scenario struct {
	Policy  Policy  `json:"policy"`
	Entries []Entry `json:"entries"`
}
type Row struct {
	Name      string   `json:"name"`
	Code      uint16   `json:"code"`
	Lifecycle string   `json:"lifecycle"`
	Zone      string   `json:"zone"`
	Issues    []string `json:"issues"`
}
type Budget struct {
	Zone     string `json:"zone"`
	Capacity uint16 `json:"capacity"`
	Occupied uint16 `json:"occupied"`
	Free     uint16 `json:"free"`
}
type Report struct {
	StableMax       uint16   `json:"stable_max"`
	ExperimentalMin uint16   `json:"experimental_min"`
	Rows            []Row    `json:"rows"`
	Budgets         []Budget `json:"budgets"`
	Clean           bool     `json:"clean"`
}

func zone(c, stable, experimental uint16) string {
	if c <= stable {
		return "stable"
	}
	if c >= experimental {
		return "experimental"
	}
	return "reserved"
}
func plan(s Scenario) (Report, error) {
	r := Report{}
	if s.Policy.StableMax == nil || s.Policy.ExperimentalMin == nil {
		return r, fmt.Errorf("policy boundaries are required")
	}
	a, b := *s.Policy.StableMax, *s.Policy.ExperimentalMin
	if a >= b || b > 255 {
		return r, fmt.Errorf("require 0 <= stable_max < experimental_min <= 255")
	}
	if len(s.Entries) > 512 {
		return r, fmt.Errorf("at most 512 entries")
	}
	names := map[string]int{}
	codes := map[uint16]int{}
	for _, e := range s.Entries {
		if e.Name == "" || e.Name != strings.TrimSpace(e.Name) || len(e.Name) > 100 || e.Code == nil || *e.Code > 255 {
			return Report{}, fmt.Errorf("entries need a trimmed 1–100 byte name and explicit code 0..255")
		}
		if e.Lifecycle != "stable" && e.Lifecycle != "experimental" {
			return Report{}, fmt.Errorf("lifecycle must be stable or experimental")
		}
		names[e.Name]++
		codes[*e.Code]++
	}
	r = Report{StableMax: a, ExperimentalMin: b, Rows: []Row{}, Clean: true, Budgets: []Budget{{Zone: "stable", Capacity: a + 1}, {Zone: "reserved", Capacity: b - a - 1}, {Zone: "experimental", Capacity: 256 - b}}}
	for _, e := range s.Entries {
		z := zone(*e.Code, a, b)
		row := Row{Name: e.Name, Code: *e.Code, Lifecycle: e.Lifecycle, Zone: z, Issues: []string{}}
		if names[e.Name] > 1 {
			row.Issues = append(row.Issues, "duplicate name")
		}
		if codes[*e.Code] > 1 {
			row.Issues = append(row.Issues, "duplicate code")
		}
		if z == "reserved" {
			row.Issues = append(row.Issues, "assignment uses reserved space")
		} else if z != e.Lifecycle {
			row.Issues = append(row.Issues, "lifecycle does not match code range")
		}
		if len(row.Issues) > 0 {
			r.Clean = false
		}
		r.Rows = append(r.Rows, row)
	}
	// Occupancy counts unique numerical codes, including policy-violating entries.
	for c := range codes {
		z := zone(c, a, b)
		for i := range r.Budgets {
			if r.Budgets[i].Zone == z {
				r.Budgets[i].Occupied++
			}
		}
	}
	for i := range r.Budgets {
		r.Budgets[i].Free = r.Budgets[i].Capacity - r.Budgets[i].Occupied
	}
	return r, nil
}

var page = template.Must(template.New("report").Funcs(template.FuncMap{"issues": func(v []string) string {
	if len(v) == 0 {
		return "none"
	}
	return strings.Join(v, "; ")
}}).Parse(`<!doctype html>
<html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width">
<title>RFC 37 · Extension registry checker</title>
<style>body{font:18px system-ui;max-width:1100px;margin:3rem auto;padding:0 1rem;background:#f4f7fb;color:#172d44}table{width:100%;border-collapse:collapse;background:white;margin:1rem 0}th,td{text-align:left;padding:.8rem;border:1px solid #ccd6df}a{color:#075fa3}</style>
<h1>RFC 37 · Extension registry checker</h1>
<p>A teaching exercise inspired by discussion about preserving space for extension and experimentation. This policy uses an 8-bit toy registry: stable codes 0–{{.StableMax}}, experimental codes {{.ExperimentalMin}}–255, and a reserved gap.</p>
<p>These boundaries are explicit exercise choices, not an RFC 37 standard allocation. No codes are allocated or registered externally.</p>
<p><strong>Registry passes supplied policy: {{.Clean}}</strong></p>
<h2>Assignments</h2><table><thead><tr><th>Name</th><th>Code</th><th>Lifecycle</th><th>Range</th><th>Issues</th></tr></thead><tbody>
{{range .Rows}}<tr><td>{{.Name}}</td><td>{{.Code}}</td><td>{{.Lifecycle}}</td><td>{{.Zone}}</td><td>{{issues .Issues}}</td></tr>{{end}}</tbody></table>
<h2>Space budget</h2><table><thead><tr><th>Range</th><th>Capacity</th><th>Occupied</th><th>Free</th></tr></thead><tbody>
{{range .Budgets}}<tr><td>{{.Zone}}</td><td>{{.Capacity}}</td><td>{{.Occupied}}</td><td>{{.Free}}</td></tr>{{end}}</tbody></table>
<p>Occupied counts distinct code values even when an assignment violates policy. Duplicate rows do not consume extra code points. Reserved free space is held aside, not available for assignment under this policy.</p>
<p>This checker validates supplied registry records, not peer negotiation, interoperability, rollout safety, or standards compliance. Historical proposals are not production recommendations.</p>
<p>Static Go-generated report, with no live backend. Run the CLI for custom input.</p>
<p><a href="https://www.rfc-editor.org/rfc/rfc37.txt">Official RFC 37</a> · <a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0024-memo-register/tree/main/cmd/registry37">Source, tests and instructions</a></p></html>
`))
