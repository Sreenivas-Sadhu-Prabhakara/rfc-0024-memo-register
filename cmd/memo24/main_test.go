package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"example.com/rfc-go-labs/lab"
)

const sample = `{"title":"<memo>","roster":["team-a"],"copies":["team-a"]}`

func TestRunJSON(t *testing.T) {
	var out bytes.Buffer
	if err := run(nil, strings.NewReader(sample), &out); err != nil {
		t.Fatal(err)
	}
	var report lab.Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Matched != 1 || len(report.Recipients) != 1 || report.Recipients[0].ID != "team-a" {
		t.Fatalf("unexpected output %+v", report)
	}
}

func TestRunHTML(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-html"}, strings.NewReader(sample), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "<h1>&lt;memo&gt;</h1>") {
		t.Fatalf("expected an escaped title in the report")
	}
}

func TestRunRejectsInvalidInput(t *testing.T) {
	for name, input := range map[string]string{
		"trailing JSON": sample + "{}",
		"unknown field": strings.Replace(sample, "title", "unknown", 1),
		"oversized":     sample + strings.Repeat(" ", 64<<10),
	} {
		t.Run(name, func(t *testing.T) {
			var out bytes.Buffer
			if err := run(nil, strings.NewReader(input), &out); err == nil || out.Len() != 0 {
				t.Fatalf("got error %v and %d output bytes", err, out.Len())
			}
		})
	}
}
