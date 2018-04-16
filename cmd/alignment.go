package cmd

import (
	"fmt"

	tw "github.com/olekukonko/tablewriter"
)

var UnknownAlignment = fmt.Errorf("unknown alignment")

type Alignment string

const (
	AlignDefault Alignment = "="
	AlignCenter Alignment = "-"
	AlignLeft Alignment = "<"
	AlignRight Alignment = ">"
	AlignEmpty Alignment = ""
)

func (a Alignment) ToString() string {
	switch a {
	case AlignDefault:
		return "="
	case AlignCenter:
		return "-"
	case AlignLeft:
		return "<"
	case AlignRight:
		return ">"
	}
	return ""
}

func (a Alignment) ToTableWriter() int {
	switch a {
	case AlignDefault:
		return tw.ALIGN_DEFAULT
	case AlignCenter:
		return tw.ALIGN_CENTER
	case AlignRight:
		return tw.ALIGN_RIGHT
	case AlignLeft:
		return tw.ALIGN_LEFT
	}
	return -1
}

func FromString(s string) (Alignment, error) {
	switch s {
	case "=":
		return AlignDefault, nil
	case "-":
		return AlignCenter, nil
	case "<":
		return AlignLeft, nil
	case ">":
		return AlignRight, nil
	}
	return AlignEmpty, UnknownAlignment
}

func FromStringArray(strs []string, out []Alignment) error {
	for _, s := range strs {
		a, err := FromString(s)
		if err != nil {
			return err
		}
		out = append(out, a)
	}
	return nil
}

