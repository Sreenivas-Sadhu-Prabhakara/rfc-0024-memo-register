package main

import (
	"fmt"
	"strings"
)

type Scenario struct {
	BudgetUS int64    `json:"budget_us"`
	Samples  []Sample `json:"samples"`
}

type Sample struct {
	Name            string `json:"name"`
	SentUS          int64  `json:"sent_us"`
	ReceivedUS      int64  `json:"received_us"`
	SenderErrorUS   int64  `json:"sender_error_us"`
	ReceiverErrorUS int64  `json:"receiver_error_us"`
}

type Result struct {
	Name       string `json:"name"`
	ObservedUS int64  `json:"observed_us"`
	ErrorUS    int64  `json:"combined_error_us"`
	MinUS      int64  `json:"min_us"`
	MaxUS      int64  `json:"max_us"`
	Consistent bool   `json:"consistent"`
	Verdict    string `json:"verdict"`
}

type Report struct {
	BudgetUS int64    `json:"budget_us"`
	Results  []Result `json:"results"`
}

// analyze assumes each endpoint's absolute timestamp error is within its supplied bound.
func analyze(input Scenario) (Report, error) {
	const limit int64 = 1_000_000_000_000
	if input.BudgetUS < 0 || input.BudgetUS > limit || len(input.Samples) == 0 {
		return Report{}, fmt.Errorf("budget_us must be 0..1e12 and samples must be nonempty")
	}
	report := Report{BudgetUS: input.BudgetUS, Results: make([]Result, 0, len(input.Samples))}
	names := make(map[string]bool)
	for _, sample := range input.Samples {
		if strings.TrimSpace(sample.Name) == "" || len(sample.Name) > 100 || names[sample.Name] {
			return Report{}, fmt.Errorf("sample names must be unique, nonblank, and at most 100 bytes")
		}
		names[sample.Name] = true
		if sample.SentUS < -limit || sample.SentUS > limit || sample.ReceivedUS < -limit || sample.ReceivedUS > limit ||
			sample.SenderErrorUS < 0 || sample.SenderErrorUS > limit || sample.ReceiverErrorUS < 0 || sample.ReceiverErrorUS > limit {
			return Report{}, fmt.Errorf("sample %q has timestamps or error bounds outside the supported range", sample.Name)
		}
		row := Result{Name: sample.Name, ObservedUS: sample.ReceivedUS - sample.SentUS,
			ErrorUS: sample.SenderErrorUS + sample.ReceiverErrorUS}
		row.MinUS = row.ObservedUS - row.ErrorUS
		row.MaxUS = row.ObservedUS + row.ErrorUS
		row.Consistent = row.MaxUS >= 0
		switch {
		case !row.Consistent:
			row.Verdict = "Inconsistent timestamps"
		default:
			row.MinUS = max(int64(0), row.MinUS)
			switch {
			case row.MaxUS <= input.BudgetUS:
				row.Verdict = "Within budget"
			case row.MinUS > input.BudgetUS:
				row.Verdict = "Exceeds budget"
			default:
				row.Verdict = "Uncertain"
			}
		}
		report.Results = append(report.Results, row)
	}
	return report, nil
}
