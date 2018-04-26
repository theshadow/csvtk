// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	tw "github.com/olekukonko/tablewriter"
)

// ErrUnknownAlignment Returned when the passed in alignment character isn't a valid Alignment type.
var ErrUnknownAlignment = fmt.Errorf("unknown alignment")

// Alignment An alignment specifies how a header, column, or footer should be aligned while rendering.
type Alignment string

const (
	// AlignDefault will align the header, column, or footer in the default manner.
	AlignDefault Alignment = "="
	// AlignCenter will align the header, column, or footer to the center.
	AlignCenter Alignment = "-"
	// AlignLeft will align the header, column, or footer to the left.
	AlignLeft Alignment = "<"
	// AlignRight will align the header, column, or footer to the right.
	AlignRight Alignment = ">"
	// AlignEmpty only used during an error state.
	AlignEmpty Alignment = ""
)

// ToString will convert the Alignment type into its associated string value.
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

// ToTableWriter will convert the alignment to a TableWriter format
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

// FromString accepts a string and converts it into an Alignment type. If the alignment is invalid it will return
// an ErrUnknownAlignment
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
	return AlignEmpty, ErrUnknownAlignment
}

// FromStringArray accepts an array of strings and attempts to convert them into an Alignment type. Any error will
// break the conversion and return an ErrUnknownAlignment
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
