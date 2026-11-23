package main

import (
	"fmt"
	"strings"
)

type ClockPair struct {
	Name                string `json:"name"`
	SourceOffsetUS      int64  `json:"source_offset_us"`
	DestinationOffsetUS int64  `json:"destination_offset_us"`
}
type Scenario struct {
	QueueUS   int64       `json:"queue_us"`
	NetworkUS int64       `json:"network_us"`
	ServiceUS int64       `json:"service_us"`
	Clocks    []ClockPair `json:"clocks"`
}
type Trace struct {
	Name                 string `json:"name"`
	ClockDifferenceUS    int64  `json:"clock_difference_us"`
	SourceStartUS        int64  `json:"source_start_us"`
	SourceSendUS         int64  `json:"source_send_us"`
	DestinationReceiveUS int64  `json:"destination_receive_us"`
	DestinationEndUS     int64  `json:"destination_end_us"`
	LocalQueueUS         int64  `json:"local_queue_us"`
	LocalServiceUS       int64  `json:"local_service_us"`
	RawTransitUS         int64  `json:"raw_transit_us"`
	RawTotalUS           int64  `json:"raw_total_us"`
}
type Report struct {
	ActualNetworkUS int64   `json:"actual_network_us"`
	ActualTotalUS   int64   `json:"actual_total_us"`
	Traces          []Trace `json:"traces"`
}

// analyze applies constant offsets to timestamps on the same synthetic request path.
func analyze(input Scenario) (Report, error) {
	const limit int64 = 1_000_000_000_000
	if input.QueueUS < 0 || input.QueueUS > limit || input.NetworkUS < 0 || input.NetworkUS > limit || input.ServiceUS < 0 || input.ServiceUS > limit || len(input.Clocks) == 0 {
		return Report{}, fmt.Errorf("durations must be 0..1e12 microseconds and clocks must be nonempty")
	}
	total := input.QueueUS + input.NetworkUS + input.ServiceUS
	report := Report{ActualNetworkUS: input.NetworkUS, ActualTotalUS: total, Traces: make([]Trace, 0, len(input.Clocks))}
	names := make(map[string]bool)
	for _, clock := range input.Clocks {
		if strings.TrimSpace(clock.Name) == "" || len(clock.Name) > 100 || names[clock.Name] {
			return Report{}, fmt.Errorf("clock names must be unique nonblank strings up to 100 bytes")
		}
		names[clock.Name] = true
		if clock.SourceOffsetUS < -limit || clock.SourceOffsetUS > limit || clock.DestinationOffsetUS < -limit || clock.DestinationOffsetUS > limit {
			return Report{}, fmt.Errorf("offsets must be within ±1e12 microseconds")
		}
		trace := Trace{
			Name: clock.Name, ClockDifferenceUS: clock.DestinationOffsetUS - clock.SourceOffsetUS,
			SourceStartUS: clock.SourceOffsetUS, SourceSendUS: input.QueueUS + clock.SourceOffsetUS,
			DestinationReceiveUS: input.QueueUS + input.NetworkUS + clock.DestinationOffsetUS,
			DestinationEndUS:     total + clock.DestinationOffsetUS,
		}
		trace.LocalQueueUS = trace.SourceSendUS - trace.SourceStartUS
		trace.LocalServiceUS = trace.DestinationEndUS - trace.DestinationReceiveUS
		trace.RawTransitUS = trace.DestinationReceiveUS - trace.SourceSendUS
		trace.RawTotalUS = trace.DestinationEndUS - trace.SourceStartUS
		report.Traces = append(report.Traces, trace)
	}
	return report, nil
}
