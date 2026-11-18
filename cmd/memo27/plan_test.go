package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestPlan(t *testing.T) {
	input := Revision{
		Number: 3,
		Roster: []string{"c", "b", "a"},
		Dispatches: []Dispatch{
			{Site: "b", Revision: 3}, {Site: "b", Revision: 1},
			{Site: "c", Revision: 2}, {Site: "c", Revision: 2},
			{Site: "removed", Revision: 3},
		},
	}
	want := Report{Revision: 3, Pending: 2, Sites: []SitePlan{
		{Site: "a", Action: "Send first copy"},
		{Site: "b", Latest: 3, Action: "Already recorded"},
		{Site: "c", Latest: 2, Action: "Send updated copy"},
	}}
	got, err := plan(input)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, err %v, want %+v", got, err, want)
	}
	if !reflect.DeepEqual(input.Roster, []string{"c", "b", "a"}) {
		t.Fatal("plan modified the roster")
	}
}

func TestRunRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"revision":0,"roster":["a"],"dispatches":[]}`,
		`{"revision":1,"roster":["a","a"],"dispatches":[]}`,
		`{"revision":1,"roster":["a"],"dispatches":[{"site":"a","revision":2}]}`,
		`{"revision":1,"roster":["Bad ID"],"dispatches":[]}`,
		`{"revision":1,"roster":["a"]}`,
		`{"revision":1,"roster":["a"],"dispatches":[]} []`,
		`{"revision":1,"roster":["a"],"dispatches":[],"typo":true}`,
		strings.Repeat(" ", (64<<10)+1),
	} {
		var output bytes.Buffer
		if err := run(nil, strings.NewReader(input), &output); err == nil || output.Len() != 0 {
			t.Errorf("expected an error and no output for input of %d bytes", len(input))
		}
	}
}

func TestRunEmptyLedger(t *testing.T) {
	const input = `{"revision":1,"roster":["new-site"],"dispatches":[]}`
	for _, test := range []struct {
		args []string
		want string
	}{
		{nil, `"pending": 1`},
		{[]string{"-html"}, "<td>None</td><td>Send first copy</td>"},
	} {
		var output bytes.Buffer
		if err := run(test.args, strings.NewReader(input), &output); err != nil || !strings.Contains(output.String(), test.want) {
			t.Errorf("missing %q, err %v", test.want, err)
		}
	}
}
