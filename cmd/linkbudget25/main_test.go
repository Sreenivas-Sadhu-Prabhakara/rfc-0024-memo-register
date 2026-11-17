package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestCompareBoundaries(t *testing.T) {
	ids := []int{0, 31, 32, 63, 64, 254, 255, -1, 256, 63}
	want := []Budget{
		{Ceiling: 31, TableSlots: 32, Data: 2, Diagnostics: 1, Unassigned: 5, Invalid: 2},
		{Ceiling: 63, TableSlots: 64, Data: 5, Diagnostics: 1, Unassigned: 2, Invalid: 2},
		{Ceiling: 254, TableSlots: 255, Data: 7, Diagnostics: 1, Unassigned: 0, Invalid: 2},
	}
	if got := compare(ids); !reflect.DeepEqual(got, want) {
		t.Fatalf("compare(%v) = %+v, want %+v", ids, got, want)
	}
}

func TestCompareEmpty(t *testing.T) {
	for _, row := range compare(nil) {
		if row.Data+row.Unassigned+row.Diagnostics+row.Invalid != 0 {
			t.Fatalf("empty trace produced counts: %+v", row)
		}
	}
}

func TestRunRejectsInvalidTrace(t *testing.T) {
	for _, input := range []string{"null", "[null]", "[1.5]", "[1] []", "{}", "[]" + strings.Repeat(" ", 64<<10)} {
		var output bytes.Buffer
		if err := run(nil, strings.NewReader(input), &output); err == nil || output.Len() != 0 {
			t.Fatal("expected an input error and no report")
		}
	}
}

func TestRunHTML(t *testing.T) {
	var output bytes.Buffer
	if err := run([]string{"-html"}, strings.NewReader("[63,64,255]"), &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "<th scope=\"row\">63</th><td>64</td><td>1</td><td>1</td><td>1</td><td>0</td>") {
		t.Fatal("missing expected row for the 64-slot policy")
	}
}
