package main

import (
	"fmt"
	"strings"
)

type Sample struct {
	Name    string `json:"name"`
	StartUS int64  `json:"start_us"`
	RTTUS   int64  `json:"rtt_us"`
}

type Scenario struct {
	Samples []Sample `json:"samples"`
}

type Measurement struct {
	Name       string `json:"name"`
	ActualUS   int64  `json:"actual_us"`
	MeasuredUS int64  `json:"measured_us"`
	ErrorUS    int64  `json:"error_us"`
}

type Clock struct {
	TickUS       int64         `json:"tick_us"`
	Measurements []Measurement `json:"measurements"`
}

type Report struct {
	Clocks []Clock `json:"clocks"`
}

// analyze floors both endpoints onto an ideal clock's tick grid.
func analyze(input Scenario) (Report, error) {
	const limit int64 = 1_000_000_000_000
	if len(input.Samples) == 0 {
		return Report{}, fmt.Errorf("samples must be nonempty")
	}
	names := make(map[string]bool)
	for _, sample := range input.Samples {
		if strings.TrimSpace(sample.Name) == "" || len(sample.Name) > 100 || names[sample.Name] {
			return Report{}, fmt.Errorf("names must be unique, nonblank, and at most 100 bytes")
		}
		names[sample.Name] = true
		if sample.StartUS < 0 || sample.StartUS > limit || sample.RTTUS < 0 || sample.RTTUS > limit {
			return Report{}, fmt.Errorf("start_us and rtt_us must be 0..1e12")
		}
	}

	report := Report{Clocks: make([]Clock, 0, 3)}
	for _, tick := range []int64{100, 1000, 25600} {
		clock := Clock{TickUS: tick, Measurements: make([]Measurement, 0, len(input.Samples))}
		for _, sample := range input.Samples {
			start := sample.StartUS / tick * tick
			end := (sample.StartUS + sample.RTTUS) / tick * tick
			measured := end - start
			clock.Measurements = append(clock.Measurements, Measurement{
				Name: sample.Name, ActualUS: sample.RTTUS,
				MeasuredUS: measured, ErrorUS: measured - sample.RTTUS,
			})
		}
		report.Clocks = append(report.Clocks, clock)
	}
	return report, nil
}
