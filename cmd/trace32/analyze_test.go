package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestOffsetCancellation(t *testing.T) {
	got, err := analyze(Scenario{QueueUS: 5, NetworkUS: 8, ServiceUS: 10, Clocks: []ClockPair{
		{Name: "aligned", SourceOffsetUS: 100, DestinationOffsetUS: 100},
		{Name: "behind", SourceOffsetUS: 50, DestinationOffsetUS: -50},
		{Name: "ahead", SourceOffsetUS: -50, DestinationOffsetUS: 50},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if got.ActualTotalUS != 23 || got.ActualNetworkUS != 8 {
		t.Fatalf("unexpected truth %+v", got)
	}
	for i, want := range []int64{8, -92, 108} {
		row := got.Traces[i]
		if row.LocalQueueUS != 5 || row.LocalServiceUS != 10 || row.RawTransitUS != want || row.RawTotalUS != want+15 {
			t.Errorf("bad trace %+v", row)
		}
	}
}
func TestZeroAndBounds(t *testing.T) {
	got, err := analyze(Scenario{Clocks: []ClockPair{{Name: "zero"}}})
	if err != nil || got.ActualTotalUS != 0 || got.Traces[0].RawTotalUS != 0 {
		t.Fatalf("zero trace %+v, %v", got, err)
	}
	got, err = analyze(Scenario{QueueUS: 1e12, NetworkUS: 1e12, ServiceUS: 1e12, Clocks: []ClockPair{{Name: "limit", SourceOffsetUS: -1e12, DestinationOffsetUS: 1e12}}})
	if err != nil || got.Traces[0].RawTotalUS != 5e12 {
		t.Fatalf("boundary trace %+v, %v", got, err)
	}
}
func TestInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"clocks":[]}`, `{"queue_us":-1,"clocks":[{"name":"x"}]}`,
		`{"network_us":1000000000001,"clocks":[{"name":"x"}]}`,
		`{"clocks":[{"name":"x","source_offset_us":-1000000000001}]}`,
		`{"clocks":[{"name":"x"},{"name":"x"}]}`,
		`{"clocks":[{"name":" "}]}`, `{"clocks":[{"name":"x"}],"typo":1}`,
		`{"clocks":[{"name":"x"}]} {}`, strings.Repeat(" ", (64<<10)+1),
	} {
		var out bytes.Buffer
		if err := run(nil, strings.NewReader(input), &out); err == nil || out.Len() != 0 {
			t.Errorf("expected rejection of %d bytes", len(input))
		}
	}
}
func TestCLIReports(t *testing.T) {
	for _, tt := range []struct {
		args []string
		want string
	}{
		{nil, `"actual_total_us": 23`}, {[]string{"-html"}, "&lt;clock&gt;"},
	} {
		var out bytes.Buffer
		err := run(tt.args, strings.NewReader(`{"queue_us":5,"network_us":8,"service_us":10,"clocks":[{"name":"<clock>"}]}`), &out)
		if err != nil || !strings.Contains(out.String(), tt.want) {
			t.Errorf("missing %q, error %v", tt.want, err)
		}
	}
}
