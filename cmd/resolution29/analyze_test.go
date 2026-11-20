package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestClockResolution(t *testing.T) {
	tests := []struct {
		name       string
		start, rtt int64
		want       [3]int64
	}{
		{"below all ticks", 0, 99, [3]int64{0, 0, 0}},
		{"crossing boundaries", 999, 2, [3]int64{100, 1000, 0}},
		{"exact tick multiples", 0, 25600, [3]int64{25600, 25000, 25600}},
		{"phase changes measurement", 25599, 1, [3]int64{100, 0, 25600}},
		{"zero duration", 25599, 0, [3]int64{0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := analyze(Scenario{Samples: []Sample{{Name: tt.name, StartUS: tt.start, RTTUS: tt.rtt}}})
			if err != nil {
				t.Fatal(err)
			}
			for i, clock := range report.Clocks {
				got := clock.Measurements[0]
				if got.MeasuredUS != tt.want[i] || got.ErrorUS != tt.want[i]-tt.rtt {
					t.Errorf("tick %d got %+v, want measured %d", clock.TickUS, got, tt.want[i])
				}
				if got.ErrorUS <= -clock.TickUS || got.ErrorUS >= clock.TickUS {
					t.Errorf("quantization error exceeds one tick: %+v", got)
				}
			}
		})
	}
}

func TestInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"samples":[]}`,
		`{"samples":[{"name":"x","start_us":-1}]}`,
		`{"samples":[{"name":"x","rtt_us":-1}]}`,
		`{"samples":[{"name":"x","rtt_us":1000000000001}]}`,
		`{"samples":[{"name":"x"},{"name":"x"}]}`,
		`{"samples":[{"name":" "}]}`,
		`{"samples":[{"name":"x"}],"typo":1}`,
		`{"samples":[{"name":"x"}]} {}`,
		strings.Repeat(" ", (64<<10)+1),
	} {
		var output bytes.Buffer
		if err := run(nil, strings.NewReader(input), &output); err == nil || output.Len() != 0 {
			t.Errorf("expected error and no output for %d bytes", len(input))
		}
	}
}

func TestHTMLAndJSON(t *testing.T) {
	const input = `{"samples":[{"name":"<clock>","start_us":900,"rtt_us":800}]}`
	for _, tt := range []struct {
		args []string
		want string
	}{
		{nil, `"measured_us": 1000`},
		{[]string{"-html"}, "<th scope=\"row\">&lt;clock&gt;</th>"},
	} {
		var output bytes.Buffer
		if err := run(tt.args, strings.NewReader(input), &output); err != nil || !strings.Contains(output.String(), tt.want) {
			t.Errorf("missing %q, error %v", tt.want, err)
		}
	}
}
