package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestPlanConsistency(t *testing.T) {
	input := Scenario{
		Roster: []Contact{{Site: "c", Name: "C"}, {Site: "b", Name: "B"}, {Site: "a", Name: "A"}},
		Directory: []Entry{
			{Site: "b", Name: "Old B", Mailbox: "b@example.invalid"},
			{Site: "a", Name: "A", Mailbox: "a@example.invalid"},
			{Site: "z", Name: "Z", Mailbox: "z@example.invalid"},
		},
	}
	got, err := plan(input)
	want := Report{Ready: 1, Unresolved: 2, Routes: []Route{
		{Site: "a", Expected: "A", Listed: "A", Mailbox: "a@example.invalid", Status: "Ready"},
		{Site: "b", Expected: "B", Listed: "Old B", Status: "Contact mismatch"},
		{Site: "c", Expected: "C", Status: "Missing directory entry"},
	}, Unused: []string{"z"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, err %v, want %+v", got, err, want)
	}
	if input.Roster[0].Site != "c" {
		t.Fatal("roster order was changed")
	}
}

func TestEmptyDirectory(t *testing.T) {
	got, err := plan(Scenario{Roster: []Contact{{Site: "a", Name: "A"}}, Directory: []Entry{}})
	if err != nil || got.Unresolved != 1 || got.Ready != 0 {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestInputValidation(t *testing.T) {
	for _, input := range []string{
		`{"roster":[],"directory":[]}`,
		`{"roster":[{"site":"a","name":"A"}]}`,
		`{"roster":[{"site":"a","name":"A"},{"site":"a","name":"B"}],"directory":[]}`,
		`{"roster":[{"site":"BAD","name":"A"}],"directory":[]}`,
		`{"roster":[{"site":"a","name":" A"}],"directory":[]}`,
		`{"roster":[{"site":"a","name":"A"}],"directory":[{"site":"a","name":"A","mailbox":"invalid"}]}`,
		`{"roster":[{"site":"a","name":"A"}],"directory":[{"site":"a","name":"A","mailbox":"A <a@example.invalid>"}]}`,
		`{"roster":[{"site":"a","name":"A"}],"directory":[{"site":"a","name":"A","mailbox":"a@example.invalid"},{"site":"a","name":"A","mailbox":"b@example.invalid"}]}`,
		`{"roster":[{"site":"a","name":"A"}],"directory":[],"typo":true}`,
		`{"roster":[{"site":"a","name":"A"}],"directory":[]} {}`,
		strings.Repeat(" ", (64<<10)+1),
	} {
		var output bytes.Buffer
		if err := run(nil, strings.NewReader(input), &output); err == nil || output.Len() != 0 {
			t.Errorf("expected error and no report for %d bytes", len(input))
		}
	}
}

func TestHTMLAndJSON(t *testing.T) {
	const input = `{"roster":[{"site":"a","name":"<A>"}],"directory":[{"site":"a","name":"<A>","mailbox":"a@example.invalid"}]}`
	for _, tt := range []struct {
		args []string
		want string
	}{
		{nil, `"ready": 1`},
		{[]string{"-html"}, "&lt;A&gt;"},
	} {
		var output bytes.Buffer
		if err := run(tt.args, strings.NewReader(input), &output); err != nil || !strings.Contains(output.String(), tt.want) {
			t.Errorf("missing %q, error %v", tt.want, err)
		}
	}
}
