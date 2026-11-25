package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestRolloverAndBounds(t *testing.T) {
	tests := []struct {
		start, end, tick, bound, count, min, max uint64
		status                                   string
	}{
		{modulus - 3, 4, 100, 100, 1, 700, 700, "unique tick count"},
		{modulus - 3, 4, 100, modulus + 7, 2, 700, 1677722300, "ambiguous"},
		{modulus - 3, 4, 100, 6, 0, 0, 0, "inconsistent"},
		{42, 42, 1000, modulus, 2, 0, 16777216000, "ambiguous"},
		{0, 0, 100, 0, 1, 0, 0, "unique tick count"},
		{0, 1, 1000, modulus, 1, 1000, 1000, "unique tick count"},
		{0, 1, 1000, modulus + 1, 2, 1000, 16777217000, "ambiguous"},
	}
	for _, v := range tests {
		r, e := plan(Scenario{[]Sample{{"case", v.start, v.end, v.tick, v.bound}}})
		if e != nil {
			t.Fatal(e)
		}
		x := r.Results[0]
		if x.CandidateCount != v.count || x.Status != v.status || x.WrapPeriodUS != modulus*v.tick {
			t.Fatalf("%+v", x)
		}
		if v.count == 0 {
			if x.MinGridUS != nil || x.MaxGridUS != nil {
				t.Fatal("fabricated duration")
			}
		} else if *x.MinGridUS != v.min || *x.MaxGridUS != v.max {
			t.Fatalf("duration %+v", x)
		}
	}
}
func TestCandidatesAgainstEnumeration(t *testing.T) {
	// Independently enumerate possible counts for fixed residues and up to three wraps.
	for _, start := range []uint64{0, 1, 42, modulus - 1} {
		for _, end := range []uint64{0, 7, 42, modulus - 1} {
			for _, bound := range []uint64{0, 6, modulus - 1, modulus, modulus + 7, 3*modulus + 100} {
				var want []uint64
				// Each integer advance has a unique decomposition into quotient and residue.
				for wraps := uint64(0); wraps <= 3; wraps++ {
					residue := (end + modulus - start) % modulus
					n := wraps*modulus + residue
					if n <= bound && (start+n)%modulus == end {
						want = append(want, n)
					}
				}
				r, e := plan(Scenario{[]Sample{{"enum", start, end, 100, bound}}})
				if e != nil {
					t.Fatal(e)
				}
				x := r.Results[0]
				if x.CandidateCount != uint64(len(want)) {
					t.Fatalf("%d %d %d: %+v", start, end, bound, x)
				}
				if len(want) > 0 && (*x.MinGridUS != want[0]*100 || *x.MaxGridUS != want[len(want)-1]*100) {
					t.Fatal("extrema")
				}
			}
		}
	}
}
func TestInvalidScenarios(t *testing.T) {
	for _, s := range []Scenario{{}, {make([]Sample, 101)}, {[]Sample{{" ", 0, 0, 100, 0}}}, {[]Sample{{"x", modulus, 0, 100, 0}}}, {[]Sample{{"x", 0, modulus, 100, 0}}}, {[]Sample{{"x", 0, 0, 10, 0}}}, {[]Sample{{"x", 0, 0, 100, maxAdvance + 1}}}, {[]Sample{{"x", 0, 0, 100, 0}, {"x", 0, 0, 100, 0}}}} {
		if _, e := plan(s); e == nil {
			t.Fatal("accepted invalid scenario")
		}
	}
	r, e := plan(Scenario{[]Sample{{"limit", 0, 0, 1000, maxAdvance}}})
	if e != nil || *r.Results[0].MaxGridUS > maxAdvance*1000 {
		t.Fatal("overflow")
	}
}
func TestCLI(t *testing.T) {
	data, e := os.ReadFile("samples.json")
	if e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e = run(nil, bytes.NewReader(data), &out); e != nil {
		t.Fatal(e)
	}
	var r Report
	if e = json.Unmarshal(out.Bytes(), &r); e != nil || len(r.Results) != 4 {
		t.Fatal("JSON", e)
	}
	out.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(data), &out); e != nil || !strings.Contains(out.String(), "inconsistent") {
		t.Fatal("HTML", e)
	}
	for _, in := range []string{"null", string(data) + "{}", strings.Repeat(" ", 65537), "{\"unknown\":true}", strings.Replace(string(data), "\"start\": 16777213", "\"start\": -1", 1)} {
		if e = run(nil, strings.NewReader(in), &bytes.Buffer{}); e == nil {
			t.Fatal("accepted invalid input")
		}
	}
}
