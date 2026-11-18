package main

import (
	"fmt"
	"regexp"
	"sort"
)

type Revision struct {
	Number     int        `json:"revision"`
	Roster     []string   `json:"roster"`
	Dispatches []Dispatch `json:"dispatches"`
}

type Dispatch struct {
	Site     string `json:"site"`
	Revision int    `json:"revision"`
}

type SitePlan struct {
	Site   string `json:"site"`
	Latest int    `json:"latest_dispatched_revision"`
	Action string `json:"action"`
}

type Report struct {
	Revision int        `json:"revision"`
	Pending  int        `json:"pending"`
	Sites    []SitePlan `json:"sites"`
}

var siteID = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func validSite(id string) bool {
	return len(id) <= 64 && siteID.MatchString(id)
}

// plan compares the current roster with recorded dispatches for one memo series.
func plan(input Revision) (Report, error) {
	if input.Number < 1 || len(input.Roster) == 0 || input.Dispatches == nil {
		return Report{}, fmt.Errorf("a positive revision, nonempty roster, and dispatches array are required")
	}

	latest := make(map[string]int, len(input.Roster))
	for _, site := range input.Roster {
		if !validSite(site) {
			return Report{}, fmt.Errorf("invalid roster site %q", site)
		}
		if _, exists := latest[site]; exists {
			return Report{}, fmt.Errorf("duplicate roster site %q", site)
		}
		latest[site] = 0
	}
	for i, dispatch := range input.Dispatches {
		if !validSite(dispatch.Site) || dispatch.Revision < 1 || dispatch.Revision > input.Number {
			return Report{}, fmt.Errorf("dispatch %d has an invalid site or revision", i+1)
		}
		if previous, exists := latest[dispatch.Site]; exists && dispatch.Revision > previous {
			latest[dispatch.Site] = dispatch.Revision
		}
	}

	sites := append([]string(nil), input.Roster...)
	sort.Strings(sites)
	report := Report{Revision: input.Number, Sites: make([]SitePlan, 0, len(sites))}
	for _, site := range sites {
		row := SitePlan{Site: site, Latest: latest[site]}
		switch row.Latest {
		case 0:
			row.Action = "Send first copy"
		case input.Number:
			row.Action = "Already recorded"
		default:
			row.Action = "Send updated copy"
		}
		if row.Latest < input.Number {
			report.Pending++
		}
		report.Sites = append(report.Sites, row)
	}
	return report, nil
}
