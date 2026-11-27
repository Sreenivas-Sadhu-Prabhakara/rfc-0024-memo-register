package main

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func runTrace(t *testing.T, events []Event) Result {
	t.Helper()
	r, e := plan(Scenario{[]Trace{{"test", events}}})
	if e != nil {
		t.Fatal(e)
	}
	return r.Results[0]
}
func TestAllCausalInterleavings(t *testing.T) {
	events := []Event{{"A", "request"}, {"B", "request"}, {"A", "receive_request"}, {"A", "ack"}, {"B", "receive_request"}, {"B", "ack"}}
	valid := 0
	var visit func([]int)
	visit = func(order []int) {
		if len(order) < 6 {
			for n := 0; n < 6; n++ {
				used := false
				for _, v := range order {
					if v == n {
						used = true
					}
				}
				if !used {
					next := append(append([]int{}, order...), n)
					visit(next)
				}
			}
			return
		}
		pos := [6]int{}
		for i, n := range order {
			pos[n] = i
		}
		if pos[0] > pos[3] || pos[0] > pos[4] || pos[1] > pos[2] || pos[1] > pos[5] {
			return
		}
		valid++
		trace := []Event{}
		for _, n := range order {
			trace = append(trace, events[n])
		}
		r := runTrace(t, trace)
		if r.Violations != 0 || !r.A.ready() || !r.B.ready() {
			t.Fatal(order, r)
		}
		observed := [6]bool{}
		for i, n := range order {
			observed[n] = true
			wantA := observed[0] && observed[2] && observed[3]
			wantB := observed[1] && observed[4] && observed[5]
			if r.Steps[i].ReadyA != wantA || r.Steps[i].ReadyB != wantB {
				t.Fatal("premature readiness", order, i)
			}
		}
	}
	visit(nil)
	if valid != 80 {
		t.Fatalf("tested %d causal orderings, want 80", valid)
	}
}
func TestPrematureSendAndLocalReadiness(t *testing.T) {
	r := runTrace(t, []Event{{"A", "request"}, {"B", "request"}, {"A", "ack"}, {"A", "try_send"}, {"A", "receive_request"}, {"A", "try_send"}, {"B", "try_send"}})
	if r.Violations != 2 || r.Steps[3].Probe != "blocked" || r.Steps[5].Probe != "allowed" || !r.Steps[5].ReadyA || r.Steps[5].ReadyB {
		t.Fatal(r)
	}
}
func TestImpossibleObservationsAndDuplicates(t *testing.T) {
	r := runTrace(t, []Event{{"A", "ack"}, {"A", "receive_request"}})
	if r.Violations != 2 || r.A != (State{}) || r.B != (State{}) {
		t.Fatal("invalid observations changed state", r)
	}
	base := []Event{{"A", "request"}, {"B", "request"}, {"A", "receive_request"}, {"A", "ack"}, {"B", "receive_request"}, {"B", "ack"}}
	r = runTrace(t, base)
	dup := runTrace(t, append(base, base...))
	if !reflect.DeepEqual(r.A, dup.A) || !reflect.DeepEqual(r.B, dup.B) || dup.Violations != 0 {
		t.Fatal("duplicate semantics")
	}
}
func TestInvalidInput(t *testing.T) {
	for _, s := range []Scenario{{}, {Traces: []Trace{{"x", nil}}}, {Traces: []Trace{{" ", []Event{{"A", "request"}}}}}, {Traces: []Trace{{"x", []Event{{"C", "request"}}}}}, {Traces: []Trace{{"x", []Event{{"A", "unknown"}}}}}, {Traces: make([]Trace, 51)}, {Traces: []Trace{{"x", make([]Event, 1001)}}}, {Traces: []Trace{{"x", []Event{{"A", "request"}}}, {"x", []Event{{"B", "request"}}}}}} {
		if _, e := plan(s); e == nil {
			t.Fatal("accepted invalid input")
		}
	}
}
func TestCLI(t *testing.T) {
	data, e := os.ReadFile("traces.json")
	if e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e = run(nil, bytes.NewReader(data), &out); e != nil {
		t.Fatal(e)
	}
	var r Report
	if e = json.Unmarshal(out.Bytes(), &r); e != nil || len(r.Results) != 2 {
		t.Fatal("JSON", e)
	}
	for _, x := range r.Results {
		if x.Violations != 1 || !x.A.ready() || !x.B.ready() {
			t.Fatal(x)
		}
	}
	out.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(data), &out); e != nil || !strings.Contains(out.String(), "allowed") {
		t.Fatal("HTML", e)
	}
	for _, s := range []string{"null", string(data) + "{}", "{\"traces\":[],\"unknown\":true}", strings.Repeat(" ", 65537)} {
		if e = run(nil, strings.NewReader(s), &bytes.Buffer{}); e == nil {
			t.Fatal("accepted invalid JSON")
		}
	}
	b, _ := json.Marshal(Scenario{[]Trace{{"<script>alert(1)</script>", []Event{{"A", "request"}}}}})
	out.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(b), &out); e != nil || strings.Contains(out.String(), "<script>") {
		t.Fatal("escaping", e)
	}
}
