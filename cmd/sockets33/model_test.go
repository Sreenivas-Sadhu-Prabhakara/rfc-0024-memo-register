package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestLabels(t *testing.T) {
	for _, tc := range []struct {
		s        Socket
		hex, dir string
	}{{Socket{65541, 1, 7}, "0100050107", "send"}, {Socket{0xffffff, 255, 255}, "ffffffffff", "send"}, {Socket{}, "0000000000", "receive"}} {
		v, e := view(tc.s)
		if e != nil || v.Hex != tc.hex || v.Direction != tc.dir {
			t.Fatalf("%+v %v", v, e)
		}
	}
	for _, s := range []Socket{{0x1000000, 0, 0}, {0, 256, 0}, {0, 0, 256}} {
		if _, e := view(s); e == nil {
			t.Fatal("accepted overflow")
		}
	}
}
func TestAllParityPairs(t *testing.T) {
	for l := uint32(0); l < 256; l++ {
		for remote := uint32(0); remote < 256; remote++ {
			p := remote
			r, e := plan(Scenario{Cases: []Case{{Name: "parity", Local: Socket{1, 2, l}, Remote: Socket{3, 4, remote}, RelabelAEN: &p}}})
			if e != nil {
				t.Fatal(e)
			}
			v := r.Cases[0]
			if v.OppositeDirections != (l%2 != remote%2) || v.Relabel.Allowed != (l%2 == remote%2) {
				t.Fatalf("parity %d %d", l, remote)
			}
			if v.Relabel.Proposed.User != 1 || v.Relabel.Proposed.Host != 2 {
				t.Fatal("changed identity")
			}
		}
	}
}
func TestValidation(t *testing.T) {
	p := uint32(256)
	for _, s := range []Scenario{{}, {Cases: make([]Case, 101)}, {Cases: []Case{{Name: " "}}}, {Cases: []Case{{Name: "x"}, {Name: "x"}}}, {Cases: []Case{{Name: "x", RelabelAEN: &p}}}, {Cases: []Case{{Name: "x", Remote: Socket{0, 256, 0}}}}} {
		if _, e := plan(s); e == nil {
			t.Fatal("accepted invalid scenario")
		}
	}
}
func TestCLI(t *testing.T) {
	data, e := os.ReadFile("sockets.json")
	if e != nil {
		t.Fatal(e)
	}
	var b bytes.Buffer
	if e = run(nil, bytes.NewReader(data), &b); e != nil {
		t.Fatal(e)
	}
	var r Report
	if e = json.Unmarshal(b.Bytes(), &r); e != nil {
		t.Fatal(e)
	}
	if len(r.Cases) != 3 || !r.Cases[0].Relabel.Allowed || r.Cases[1].OppositeDirections || r.Cases[2].Relabel.Allowed {
		t.Fatal("fixture outcomes")
	}
	b.Reset()
	if e = run([]string{"-html"}, bytes.NewReader(data), &b); e != nil || !strings.Contains(b.String(), "0000000258") {
		t.Fatal("HTML", e)
	}
	for _, input := range []string{string(data) + "{}", strings.Replace(string(data), "\"user\": 0", "\"user\": -1", 1), "{\"cases\":[],\"unknown\":true}", strings.Repeat(" ", 65537), "null"} {
		if e = run(nil, strings.NewReader(input), &bytes.Buffer{}); e == nil {
			t.Fatal("accepted invalid input")
		}
	}
	zero := uint32(0)
	r, e = plan(Scenario{Cases: []Case{{Name: "zero", RelabelAEN: &zero}}})
	if e != nil || r.Cases[0].Relabel == nil {
		t.Fatal("lost explicit zero")
	}
}
