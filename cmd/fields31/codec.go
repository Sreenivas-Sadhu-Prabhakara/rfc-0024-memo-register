package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

type Field struct {
	Name  string `json:"name"`
	Bits  int    `json:"bits"`
	Value uint64 `json:"value"`
}
type Scenario struct {
	Fields []Field `json:"fields"`
}
type FieldView struct {
	Name    string `json:"name"`
	Offset  int    `json:"offset"`
	Bits    int    `json:"bits"`
	Value   uint64 `json:"value"`
	Decoded uint64 `json:"decoded"`
	Stream  string `json:"stream"`
}
type Report struct {
	BitCount int         `json:"bit_count"`
	Hex      string      `json:"hex"`
	Stream   string      `json:"stream"`
	Fields   []FieldView `json:"fields"`
}

var fieldName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,31}$`)

func layoutBits(fields []Field) (int, error) {
	if len(fields) == 0 {
		return 0, fmt.Errorf("fields must be nonempty")
	}
	names := make(map[string]bool)
	total := 0
	for _, f := range fields {
		if !fieldName.MatchString(f.Name) || names[f.Name] || f.Bits < 1 || f.Bits > 32 {
			return 0, fmt.Errorf("invalid or duplicate field %q, or width outside 1..32", f.Name)
		}
		names[f.Name] = true
		total += f.Bits
		if total > 512 {
			return 0, fmt.Errorf("layout exceeds 512 bits")
		}
	}
	return total, nil
}

// decode reads each field least-significant bit first from consecutive stream positions.
func decode(data []byte, fields []Field) ([]uint64, error) {
	total, err := layoutBits(fields)
	if err != nil {
		return nil, err
	}
	if len(data) != (total+7)/8 {
		return nil, fmt.Errorf("payload length does not match layout")
	}
	if remaining := total % 8; remaining != 0 && data[len(data)-1]>>remaining != 0 {
		return nil, fmt.Errorf("nonzero padding bits")
	}
	values := make([]uint64, len(fields))
	position := 0
	for i, f := range fields {
		for bit := 0; bit < f.Bits; bit++ {
			values[i] |= uint64((data[position/8]>>uint(position%8))&1) << uint(bit)
			position++
		}
	}
	return values, nil
}

func plan(input Scenario) (Report, error) {
	total, err := layoutBits(input.Fields)
	if err != nil {
		return Report{}, err
	}
	for _, f := range input.Fields {
		if f.Value >= uint64(1)<<uint(f.Bits) {
			return Report{}, fmt.Errorf("value does not fit field %q", f.Name)
		}
	}
	data := make([]byte, (total+7)/8)
	var stream strings.Builder
	views := make([]FieldView, 0, len(input.Fields))
	position := 0
	for _, f := range input.Fields {
		start := position
		for bit := 0; bit < f.Bits; bit++ {
			v := byte((f.Value >> uint(bit)) & 1)
			data[position/8] |= v << uint(position%8)
			stream.WriteByte('0' + v)
			position++
		}
		views = append(views, FieldView{Name: f.Name, Offset: start, Bits: f.Bits, Value: f.Value})
	}
	decoded, err := decode(data, input.Fields)
	if err != nil {
		return Report{}, err
	}
	bits := stream.String()
	for i := range views {
		views[i].Decoded = decoded[i]
		views[i].Stream = bits[views[i].Offset : views[i].Offset+views[i].Bits]
	}
	return Report{BitCount: total, Hex: hex.EncodeToString(data), Stream: bits, Fields: views}, nil
}
