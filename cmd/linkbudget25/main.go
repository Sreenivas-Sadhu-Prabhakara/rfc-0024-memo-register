package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
)

type Budget struct {
	Ceiling     int `json:"ceiling"`
	TableSlots  int `json:"table_slots"`
	Data        int `json:"data"`
	Diagnostics int `json:"diagnostics"`
	Unassigned  int `json:"unassigned"`
	Invalid     int `json:"invalid"`
}

func compare(ids []int) []Budget {
	budgets := make([]Budget, 0, 3)
	for _, ceiling := range []int{31, 63, 254} {
		result := Budget{Ceiling: ceiling, TableSlots: ceiling + 1}
		for _, id := range ids {
			switch {
			case id < 0 || id > 255:
				result.Invalid++
			case id == 255:
				result.Diagnostics++
			case id <= ceiling:
				result.Data++
			default:
				result.Unassigned++
			}
		}
		budgets = append(budgets, result)
	}
	return budgets
}

//go:embed report.html
var reportHTML string

var reportPage = template.Must(template.New("report").Parse(reportHTML))

func run(args []string, input io.Reader, output io.Writer) error {
	flags := flag.NewFlagSet("linkbudget25", flag.ContinueOnError)
	html := flags.Bool("html", false, "render an HTML report")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("usage: linkbudget25 [-html] < trace.json")
	}
	const maxInput = 64 << 10
	data, err := io.ReadAll(io.LimitReader(input, maxInput+1))
	if err != nil {
		return err
	}
	if len(data) > maxInput {
		return fmt.Errorf("trace exceeds 64 KiB")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	var values []*int
	if err := decoder.Decode(&values); err != nil || values == nil {
		return fmt.Errorf("trace must be a JSON array of integers")
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != io.EOF {
		return fmt.Errorf("expected one JSON array")
	}
	ids := make([]int, 0, len(values))
	for index, value := range values {
		if value == nil {
			return fmt.Errorf("trace entry %d must be an integer", index)
		}
		ids = append(ids, *value)
	}

	result := compare(ids)
	if *html {
		return reportPage.Execute(output, result)
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil && !errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
