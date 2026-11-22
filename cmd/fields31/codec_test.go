package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestGoldenEncoding(t *testing.T) {
	fields := []Field{{Name: "A", Bits: 3, Value: 5}, {Name: "B", Bits: 5, Value: 17}, {Name: "C", Bits: 9, Value: 257}}
	got, err := plan(Scenario{Fields: fields})
	if err != nil || got.Hex != "8d0101" || got.Stream != "10110001100000001" || got.BitCount != 17 {
		t.Fatalf("got %+v, error %v", got, err)
	}
	want, err := decode([]byte{0x8d, 0x01, 0x01}, fields)
	if err != nil || !reflect.DeepEqual(want, []uint64{5, 17, 257}) {
		t.Fatalf("got %v, error %v", want, err)
	}
}

func TestMaximumAndZero(t *testing.T) {
	fields := []Field{{Name: "MAX", Bits: 32, Value: 4294967295}, {Name: "ZERO", Bits: 1}}
	got, err := plan(Scenario{Fields: fields})
	if err != nil || got.Hex != "ffffffff00" || got.Fields[0].Decoded != 4294967295 || got.Fields[1].Decoded != 0 {
		t.Fatalf("got %+v, error %v", got, err)
	}
}

func TestDecodeRejectsBadFraming(t *testing.T) {
	fields := []Field{{Name: "A", Bits: 9}}
	for _, payload := range [][]byte{{1}, {1, 0, 0}, {1, 2}} {
		if _, err := decode(payload, fields); err == nil {
			t.Errorf("accepted %x", payload)
		}
	}
}

func TestInvalidInputs(t *testing.T) {
	for _, input := range []string{
		`{"fields":[]}`,
		`{"fields":[{"name":"A","bits":0}]}`,
		`{"fields":[{"name":"A","bits":33}]}`,
		`{"fields":[{"name":"A","bits":2,"value":4}]}`,
		`{"fields":[{"name":"A","bits":2,"value":-1}]}`,
		`{"fields":[{"name":"A","bits":2},{"name":"A","bits":2}]}`,
		`{"fields":[{"name":"bad name","bits":1}]}`,
		`{"fields":[{"name":"A","bits":1}],"typo":1}`,
		`{"fields":[{"name":"A","bits":1}]} {}`,
		strings.Repeat(" ", (64<<10)+1),
	} {
		var out bytes.Buffer
		if err := run(nil, strings.NewReader(input), &out); err == nil || out.Len() != 0 {
			t.Errorf("expected error and no report for %d bytes", len(input))
		}
	}
	fields := make([]Field, 17)
	for i := range fields {
		fields[i] = Field{Name: string(rune('A' + i)), Bits: 32}
	}
	if _, err := plan(Scenario{Fields: fields}); err == nil {
		t.Fatal("accepted layout over 512 bits")
	}
}

func TestCLIReports(t *testing.T) {
	for _, tt := range []struct {
		args []string
		want string
	}{
		{nil, `"hex": "05"`}, {[]string{"-html"}, "<code>101</code>"},
	} {
		var out bytes.Buffer
		err := run(tt.args, strings.NewReader(`{"fields":[{"name":"A","bits":3,"value":5}]}`), &out)
		if err != nil || !strings.Contains(out.String(), tt.want) {
			t.Errorf("missing %q, error %v", tt.want, err)
		}
	}
}
