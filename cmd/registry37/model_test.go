package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func ptr(v uint16) *uint16               { return &v }
func scenario(entries ...Entry) Scenario { return Scenario{Policy{ptr(31), ptr(224)}, entries} }
func inspect(t *testing.T, s Scenario) Report {
	t.Helper()
	r, e := plan(s)
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func TestEveryCodeAndLifecycle(t *testing.T) {
	for code := uint16(0); code < 256; code++ {
		for _, life := range []string{"stable", "experimental"} {
			r := inspect(t, scenario(Entry{"x", ptr(code), life}))
			expected := (life == "stable" && code <= 31) || (life == "experimental" && code >= 224)
			if r.Clean != expected {
				t.Fatalf("%d %s: %+v", code, life, r)
			}
			var cap, occupied, free uint16
			for _, b := range r.Budgets {
				cap += b.Capacity
				occupied += b.Occupied
				free += b.Free
			}
			if cap != 256 || occupied != 1 || free != 255 {
				t.Fatal("budget", r.Budgets)
			}
		}
	}
}
func TestCollisionsAndBudgets(t *testing.T) {
	r := inspect(t, scenario(Entry{"same", ptr(1), "stable"}, Entry{"same", ptr(1), "stable"}, Entry{"gap", ptr(64), "experimental"}, Entry{"trial", ptr(255), "experimental"}))
	if r.Clean || len(r.Rows[0].Issues) != 2 || len(r.Rows[1].Issues) != 2 || len(r.Rows[2].Issues) != 1 {
		t.Fatal(r)
	}
	for _, b := range r.Budgets {
		expected := uint16(32)
		if b.Zone == "reserved" {
			expected = 192
		}
		if b.Capacity != expected || b.Occupied != 1 || b.Free != expected-1 {
			t.Fatal(b)
		}
	}
}
func TestPoliciesAndEmptyRegistry(t *testing.T) {
	for _, p := range []Policy{{ptr(0), ptr(1)}, {ptr(254), ptr(255)}, {ptr(31), ptr(224)}} {
		r := inspect(t, Scenario{Policy: p})
		if !r.Clean || len(r.Rows) != 0 {
			t.Fatal(r)
		}
		var cap uint16
		for _, b := range r.Budgets {
			cap += b.Capacity
			if b.Free != b.Capacity {
				t.Fatal(b)
			}
		}
		if cap != 256 {
			t.Fatal(cap)
		}
	}
	r := inspect(t, scenario(Entry{"zero", ptr(0), "stable"}))
	if !r.Clean {
		t.Fatal("explicit zero")
	}
}
func TestInvalidInput(t *testing.T) {
	for _, s := range []Scenario{
		{}, {Policy: Policy{nil, ptr(224)}}, {Policy: Policy{ptr(32), ptr(32)}}, {Policy: Policy{ptr(31), ptr(256)}},
		scenario(Entry{"x", nil, "stable"}), scenario(Entry{"x", ptr(256), "stable"}), scenario(Entry{" x", ptr(1), "stable"}),
		scenario(Entry{"x", ptr(1), "unknown"}), scenario(Entry{strings.Repeat("x", 101), ptr(1), "stable"}),
		{Policy: Policy{ptr(31), ptr(224)}, Entries: make([]Entry, 513)},
	} {
		if _, e := plan(s); e == nil {
			t.Fatal("accepted invalid input")
		}
	}
}
func TestCLI(t *testing.T) {
	data, e := os.ReadFile("registry.json")
	if e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e = run(nil, bytes.NewReader(data), &out); e != nil {
		t.Fatal(e)
	}
	var r Report
	if e = json.Unmarshal(out.Bytes(), &r); e != nil || r.Clean || len(r.Rows) != 5 {
		t.Fatal("JSON", e, r)
	}
	out.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(data), &out); e != nil || !strings.Contains(out.String(), "duplicate code") {
		t.Fatal("HTML", e)
	}
	for _, s := range []string{"null", string(data) + "{}", "{\"policy\":{},\"unknown\":true}", strings.Repeat(" ", 65537), strings.Replace(string(data), "\"code\": 1", "\"code\": -1", 1)} {
		if e = run(nil, strings.NewReader(s), &bytes.Buffer{}); e == nil {
			t.Fatal("accepted invalid JSON")
		}
	}
	b, _ := json.Marshal(scenario(Entry{"<script>alert(1)</script>", ptr(1), "stable"}))
	out.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(b), &out); e != nil || strings.Contains(out.String(), "<script>") {
		t.Fatal("escaping", e)
	}
}
