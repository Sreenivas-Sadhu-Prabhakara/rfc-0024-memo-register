package main

import (
	"fmt"
	"net/mail"
	"regexp"
	"sort"
	"strings"
)

type Contact struct {
	Site string `json:"site"`
	Name string `json:"name"`
}

type Entry struct {
	Site    string `json:"site"`
	Name    string `json:"name"`
	Mailbox string `json:"mailbox"`
}

type Scenario struct {
	Roster    []Contact `json:"roster"`
	Directory []Entry   `json:"directory"`
}

type Route struct {
	Site     string `json:"site"`
	Expected string `json:"expected_contact"`
	Listed   string `json:"listed_contact"`
	Mailbox  string `json:"ready_mailbox,omitempty"`
	Status   string `json:"status"`
}

type Report struct {
	Ready      int      `json:"ready"`
	Unresolved int      `json:"unresolved"`
	Routes     []Route  `json:"routes"`
	Unused     []string `json:"unused_sites"`
}

var siteID = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func validSite(id string) bool {
	return len(id) <= 64 && siteID.MatchString(id)
}

func validName(name string) bool {
	return strings.TrimSpace(name) != "" && strings.TrimSpace(name) == name && len(name) <= 100
}

// plan joins a current distribution roster with a separately maintained directory.
func plan(input Scenario) (Report, error) {
	if len(input.Roster) == 0 || input.Directory == nil {
		return Report{}, fmt.Errorf("a nonempty roster and a directory array are required")
	}
	roster := make(map[string]string, len(input.Roster))
	for _, contact := range input.Roster {
		if !validSite(contact.Site) || !validName(contact.Name) {
			return Report{}, fmt.Errorf("invalid roster contact")
		}
		if _, exists := roster[contact.Site]; exists {
			return Report{}, fmt.Errorf("duplicate roster site %q", contact.Site)
		}
		roster[contact.Site] = contact.Name
	}
	directory := make(map[string]Entry, len(input.Directory))
	for _, entry := range input.Directory {
		address, err := mail.ParseAddress(entry.Mailbox)
		if !validSite(entry.Site) || !validName(entry.Name) || err != nil ||
			len(entry.Mailbox) > 254 || address.Address != entry.Mailbox {
			return Report{}, fmt.Errorf("invalid directory entry for %q", entry.Site)
		}
		if _, exists := directory[entry.Site]; exists {
			return Report{}, fmt.Errorf("duplicate directory site %q", entry.Site)
		}
		directory[entry.Site] = entry
	}

	sites := make([]string, 0, len(roster))
	for site := range roster {
		sites = append(sites, site)
	}
	sort.Strings(sites)
	report := Report{Routes: make([]Route, 0, len(sites)), Unused: []string{}}
	for _, site := range sites {
		entry, exists := directory[site]
		row := Route{Site: site, Expected: roster[site], Listed: entry.Name}
		switch {
		case !exists:
			row.Status = "Missing directory entry"
		case entry.Name != row.Expected:
			row.Status = "Contact mismatch"
		default:
			row.Status = "Ready"
			row.Mailbox = entry.Mailbox
			report.Ready++
		}
		report.Routes = append(report.Routes, row)
	}
	report.Unresolved = len(sites) - report.Ready
	for site := range directory {
		if _, exists := roster[site]; !exists {
			report.Unused = append(report.Unused, site)
		}
	}
	sort.Strings(report.Unused)
	return report, nil
}
