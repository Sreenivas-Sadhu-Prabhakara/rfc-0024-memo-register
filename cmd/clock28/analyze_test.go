package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestDelayIntervals(t *testing.T) {
	tests := []struct {
		name                                      string
		delta, senderError, receiverError, budget int64
		min, max                                  int64
		consistent                                bool
		verdict                                   string
	}{
		{"exact", 8, 0, 0, 10, 8, 8, true, "Within budget"},
		{"boundary", 8, 1, 1, 10, 6, 10, true, "Within budget"},
		{"overlap", 8, 6, 6, 10, 0, 20, true, "Uncertain"},
		{"slow", 15, 1, 1, 10, 13, 17, true, "Exceeds budget"},
		{"negative allowed", -2, 3, 0, 10, 0, 1, true, "Within budget"},
		{"inconsistent", -2, 0, 0, 10, -2, -2, false, "Inconsistent timestamps"},
		{"zero", 0, 0, 0, 0, 0, 0, true, "Within budget"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := analyze(Scenario{BudgetUS: tt.budget, Samples: []Sample{{
				Name: tt.name, ReceivedUS: tt.delta, SenderErrorUS: tt.senderError, ReceiverErrorUS: tt.receiverError,
			}}})
			want := Result{Name: tt.name, ObservedUS: tt.delta, ErrorUS: tt.senderError + tt.receiverError,
				MinUS: tt.min, MaxUS: tt.max, Consistent: tt.consistent, Verdict: tt.verdict}
			if err != nil || !reflect.DeepEqual(got.Results, []Result{want}) {
				t.Fatalf("got %+v, error %v, want %+v", got, err, want)
			}
		})
	}
}

func TestRunInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"budget_us":-1,"samples":[{"name":"x"}]}`,
		`{"samples":[]}`,
		`{"samples":[{"name":"x","sender_error_us":-1}]}`,
		`{"samples":[{"name":"x","received_us":1000000000001}]}`,
		`{"samples":[{"name":"x"},{"name":"x"}]}`,
		`{"samples":[{"name":" "}]}`,
		`{"samples":[{"name":"x"}],"typo":1}`,
		`{"samples":[{"name":"x"}]} {}`,
		strings.Repeat(" ", (64<<10)+1),
	} {
		var output bytes.Buffer
		if err := run(nil, strings.NewReader(input), &output); err == nil || output.Len() != 0 {
			t.Errorf("expected error and no output for %d-byte input", len(input))
		}
	}
}

func TestHTMLAndJSON(t *testing.T) {
	const input = `{"budget_us":10,"samples":[{"name":"<clock>","received_us":8,"sender_error_us":1,"receiver_error_us":1}]}`
	for _, tt := range []struct {
		args []string
		want string
	}{
		{nil, `"verdict": "Within budget"`},
		{[]string{"-html"}, "<th scope=\"row\">&lt;clock&gt;</th>"},
	} {
		var output bytes.Buffer
		if err := run(tt.args, strings.NewReader(input), &output); err != nil || !strings.Contains(output.String(), tt.want) {
			t.Errorf("missing %q, error %v", tt.want, err)
		}
	}
}
