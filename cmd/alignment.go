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

