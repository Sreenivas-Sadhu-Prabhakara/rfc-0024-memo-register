package lab

import (
	"reflect"
	"testing"
)

func TestAuditFindings(t *testing.T) {
	got, err := Audit(Register{
		Title: "Cutover", Roster: []string{"b", "a", "c"},
		Copies: []string{"z", "a", "d", "b", "z", "a"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := Report{
		Title: "Cutover", Matched: 1, Missing: 1, Duplicate: 1, Unexpected: 2,
		Recipients: []Recipient{
			{ID: "b", Expected: 1, Recorded: 1, Status: "Matched"},
			{ID: "a", Expected: 1, Recorded: 2, Status: "Duplicate"},
			{ID: "c", Expected: 1, Status: "Missing"},
			{ID: "d", Recorded: 1, Status: "Unexpected"},
			{ID: "z", Recorded: 2, Status: "Unexpected"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestAuditEmptyLedger(t *testing.T) {
	got, err := Audit(Register{Title: "Memo", Roster: []string{"a", "b"}, Copies: []string{}})
	if err != nil || got.Missing != 2 {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestAuditInvalidRegister(t *testing.T) {
	for _, register := range []Register{
		{Title: " ", Roster: []string{"a"}, Copies: []string{}},
		{Title: "Memo", Copies: []string{}},
		{Title: "Memo", Roster: []string{"a"}},
		{Title: "Memo", Roster: []string{"a", "a"}, Copies: []string{}},
		{Title: "Memo", Roster: []string{"A"}, Copies: []string{}},
		{Title: "Memo", Roster: []string{"a"}, Copies: []string{"a--b"}},
	} {
		if _, err := Audit(register); err == nil {
			t.Errorf("accepted invalid register %+v", register)
		}
	}
}
