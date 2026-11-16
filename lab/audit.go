package lab

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Register struct {
	Title  string   `json:"title"`
	Roster []string `json:"roster"`
	Copies []string `json:"copies"`
}

type Recipient struct {
	ID       string `json:"id"`
	Expected int    `json:"expected"`
	Recorded int    `json:"recorded"`
	Status   string `json:"status"`
}

type Report struct {
	Title      string      `json:"title"`
	Matched    int         `json:"matched"`
	Missing    int         `json:"missing"`
	Duplicate  int         `json:"duplicate"`
	Unexpected int         `json:"unexpected"`
	Recipients []Recipient `json:"recipients"`
}

var recipientID = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Audit compares recorded dispatches with the rule of one copy per roster member.
func Audit(register Register) (Report, error) {
	if strings.TrimSpace(register.Title) == "" || len(register.Title) > 200 {
		return Report{}, fmt.Errorf("title must contain text and be at most 200 bytes")
	}
	if len(register.Roster) == 0 || register.Copies == nil {
		return Report{}, fmt.Errorf("a nonempty roster and a copies array are required")
	}
	expected := make(map[string]bool, len(register.Roster))
	for _, id := range register.Roster {
		if len(id) > 64 || !recipientID.MatchString(id) {
			return Report{}, fmt.Errorf("invalid roster recipient %q", id)
		}
		if expected[id] {
			return Report{}, fmt.Errorf("repeated roster recipient %q", id)
		}
		expected[id] = true
	}
	counts := make(map[string]int)
	for _, id := range register.Copies {
		if len(id) > 64 || !recipientID.MatchString(id) {
			return Report{}, fmt.Errorf("invalid copy recipient %q", id)
		}
		counts[id]++
	}

	report := Report{Title: register.Title}
	for _, id := range register.Roster {
		row := Recipient{ID: id, Expected: 1, Recorded: counts[id]}
		switch row.Recorded {
		case 0:
			row.Status = "Missing"
			report.Missing++
		case 1:
			row.Status = "Matched"
			report.Matched++
		default:
			row.Status = "Duplicate"
			report.Duplicate++
		}
		report.Recipients = append(report.Recipients, row)
	}
	var unexpected []string
	for id := range counts {
		if !expected[id] {
			unexpected = append(unexpected, id)
		}
	}
	sort.Strings(unexpected)
	for _, id := range unexpected {
		report.Unexpected++
		report.Recipients = append(report.Recipients, Recipient{
			ID: id, Recorded: counts[id], Status: "Unexpected",
		})
	}
	return report, nil
}
