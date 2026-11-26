package main

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func good() Meeting {
	return Meeting{Name: "review", DraftRevision: "v3", Sites: []string{"A", "B"}, Delegates: []Delegate{{"A", "alice", "v3"}, {"B", "bob", "v3"}}}
}
func check(t *testing.T, m Meeting) Result {
	t.Helper()
	r, e := plan(Scenario{[]Meeting{m}})
	if e != nil {
		t.Fatal(e)
	}
	return r.Results[0]
}
func TestGates(t *testing.T) {
	if r := check(t, good()); r.Decision != "ready to circulate" {
		t.Fatal(r)
	}
	for _, tc := range []struct {
		name   string
		change func(*Meeting)
		field  func(Result) []string
		want   []string
	}{
		{"missing", func(m *Meeting) { m.Delegates = m.Delegates[:1] }, func(r Result) []string { return r.Missing }, []string{"B"}},
		{"multiple", func(m *Meeting) { m.Delegates = append(m.Delegates, Delegate{"A", "other", "v3"}) }, func(r Result) []string { return r.Multiple }, []string{"A"}},
		{"unexpected", func(m *Meeting) { m.Delegates = append(m.Delegates, Delegate{"X", "x", "v3"}) }, func(r Result) []string { return r.Unexpected }, []string{"X"}},
		{"stale", func(m *Meeting) { m.Delegates[0].ReviewedRevision = "v2" }, func(r Result) []string { return r.Stale }, []string{"A"}},
		{"unreviewed", func(m *Meeting) { m.Delegates[0].ReviewedRevision = "" }, func(r Result) []string { return r.Stale }, []string{"A"}},
		{"blocker", func(m *Meeting) { m.Blockers = []string{"retries"} }, func(r Result) []string { return r.Blockers }, []string{"retries"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := good()
			tc.change(&m)
			r := check(t, m)
			if r.Decision != "hold" || !reflect.DeepEqual(tc.field(r), tc.want) {
				t.Fatal(r)
			}
		})
	}
}
func TestDeterminism(t *testing.T) {
	m := good()
	m.Sites = []string{"Z", "C", "A", "B"}
	m.Delegates = append(m.Delegates, Delegate{"Y", "y", "old"}, Delegate{"X", "x", "old"})
	r := check(t, m)
	if !reflect.DeepEqual(r.Missing, []string{"C", "Z"}) || !reflect.DeepEqual(r.Unexpected, []string{"X", "Y"}) || !reflect.DeepEqual(r.Stale, []string{"X", "Y"}) {
		t.Fatal(r)
	}
	m.Sites = []string{"B", "A", "C", "Z"}
	m.Delegates[0], m.Delegates[1] = m.Delegates[1], m.Delegates[0]
	if !reflect.DeepEqual(r, check(t, m)) {
		t.Fatal("order changes result")
	}
}
func TestInvalidInput(t *testing.T) {
	if _, e := plan(Scenario{}); e == nil {
		t.Fatal("empty scenario")
	}
	for _, change := range []func(*Meeting){
		func(m *Meeting) { m.Sites = nil }, func(m *Meeting) { m.Sites = []string{"A", "A"} }, func(m *Meeting) { m.Sites = []string{" A"} },
		func(m *Meeting) { m.DraftRevision = "" }, func(m *Meeting) { m.Delegates[0].Name = "" }, func(m *Meeting) { m.Blockers = []string{" "} },
		func(m *Meeting) { m.Name = strings.Repeat("x", 101) }, func(m *Meeting) { m.Delegates = make([]Delegate, 201) },
	} {
		m := good()
		change(&m)
		if _, e := plan(Scenario{[]Meeting{m}}); e == nil {
			t.Fatal("accepted invalid input")
		}
	}
	m := good()
	if _, e := plan(Scenario{[]Meeting{m, m}}); e == nil {
		t.Fatal("duplicate meeting")
	}
}
func TestCLI(t *testing.T) {
	data, e := os.ReadFile("reviews.json")
	if e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e = run(nil, bytes.NewReader(data), &out); e != nil {
		t.Fatal(e)
	}
	var r Report
	if e = json.Unmarshal(out.Bytes(), &r); e != nil || len(r.Results) != 3 || r.Results[0].Decision != "ready to circulate" || r.Results[1].Decision != "hold" || r.Results[2].Decision != "hold" {
		t.Fatal(r, e)
	}
	out.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(data), &out); e != nil || !strings.Contains(out.String(), "Retry semantics unresolved") {
		t.Fatal("HTML", e)
	}
	for _, s := range []string{"null", string(data) + "{}", "{\"meetings\":[],\"x\":1}", strings.Repeat(" ", 65537)} {
		if e = run(nil, strings.NewReader(s), &bytes.Buffer{}); e == nil {
			t.Fatal("accepted invalid JSON")
		}
	}
	m := good()
	m.Name = "<script>alert(1)</script>"
	b, _ := json.Marshal(Scenario{[]Meeting{m}})
	out.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(b), &out); e != nil || strings.Contains(out.String(), "<script>") {
		t.Fatal("HTML escaping", e)
	}
}
