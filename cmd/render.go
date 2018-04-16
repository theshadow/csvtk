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
	"encoding/csv"
	"fmt"
	"os"
	"io"

	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type RenderOptions struct {
	Align Alignment
	AlignColumns []Alignment
	AlignHeader Alignment
	AlignFooter Alignment
	AutoFormattingHeaders bool
	AutoMergeCells bool
	AutoWrap bool
	Caption string
	CenterSeparator string
	ColumnSeparator string
	ColWidth int
	FirstRowAsHeader bool
	Footer []string
	Header []string
	Newline string
	Reflow bool
	RowLine bool
	RowSeparator string
}

var renderOpts RenderOptions

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Accepts data in the CSV format and renders it in an ASCII table format",
	Long: fmt.Sprintf(`Accepts data in the CSV format and renders it in an ASCII table format

Input:
  By default render reads from Stdin. You may use the -i|--input flag to specify an input file.

Output:
  By default render writes to Stdout. You may use the -o|--output flag to specify an output file.

Alignment:
  Alignment is specified using the '%s', '%s', '%s', and '%s' characters for left, right, center, and default respectively.`,
		AlignLeft, AlignRight, AlignCenter, AlignDefault),
	Example: `  Replacing the header in a CSV stream with an inline header.
    - tail -n +1 contact-list.csv | csv render --header "FName,LName,Email,Phone"

  Using the first line in the CSV stream as an inline header.
    - cat contact-list.csv | csv render --first-row-as-header
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Render(os.Stdin, os.Stdout, renderOpts)
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
	align := *renderCmd.Flags().StringP("align", "a", string(AlignDefault),
		"Defines the text alignment for the entire table.")
	alignCols := *renderCmd.Flags().StringArrayP("align-columns", "", []string{},
	"Defines the alignment for each column individually.")
	alignHeader := *renderCmd.Flags().StringP("align-header", "", string(AlignDefault),
		"Defines the alignment for the header.")
	alignFooter := *renderCmd.Flags().StringP("align-footer", "", string(AlignDefault),
		"Defines the alignment for the header columns.")
	autoFormattingHeaders := *renderCmd.Flags().BoolP("auto-formatting-headers", "", true,
		"When specified auto formatting of the headers will be disabled.")
	autoMergeCells := *renderCmd.Flags().BoolP("auto-merge-cells", "", false,
		"Defines the alignment for each footer column individually.")
	autoWrap := *renderCmd.Flags().BoolP("auto-wrap", "", true,
		"When set the text will not be automatically wrapped.")
	caption := *renderCmd.Flags().StringP("caption", "c", "",
		"Defines the caption message to be displayed with the table.")
	centerSeparator := *renderCmd.Flags().StringP("center-separator", "", tw.CENTER,
		"Defines what character will separate center columns.")
	columnSeparator := *renderCmd.Flags().StringP("column-separator", "", tw.COLUMN,
		"Defines what character will separate each column.")
	colWidth := *renderCmd.Flags().IntP("col-width", "w", tw.MAX_ROW_WIDTH,
		"Defines the fixed width for each column.")
	firstRowAsHeader := *renderCmd.Flags().BoolP("first-row-as-header", "", false,
		"When specified the first row will be treated as the headers.")
	footer := *renderCmd.Flags().StringArrayP("footer", "", []string{},
	"Defines what the footer columns should be.")
	header := *renderCmd.Flags().StringArrayP("header", "", []string{},
	"Defines what the header columns should be.")
	/*input := *renderCmd.Flags().StringP("input", "i", "",
		"Define a file to read from instead of Stdin.")*/
	newline := *renderCmd.Flags().StringP("newline", "", tw.NEWLINE,
		"Defines what character will be used at the end of a line.")
	reflow := *renderCmd.Flags().BoolP("reflow", "", true,
		"When specified the text will not be automatically re-flowed.")
	/*output := *renderCmd.Flags().StringP("output", "o", "",
		"Define a file to write to instead of Stdout.")*/
	rowLine := *renderCmd.Flags().BoolP("row-line", "", false,
		"When specified each row will be delimited with a row line.")
	rowSeparator := *renderCmd.Flags().StringP("row-separator", "", tw.ROW,
		"Defines the character used to separate each row.")

	aln, err := FromString(align)
	if err != nil {
	}

	var alnCols []Alignment
	err = FromStringArray(alignCols, alnCols)
	if err != nil {
	}

	var alnHeader Alignment
	alnHeader, err = FromString(alignHeader)

	var alnFooter Alignment
	alnFooter, err = FromString(alignFooter)

	renderOpts = RenderOptions{
		Align: aln,
		AlignColumns: alnCols,
		AlignHeader: alnHeader,
		AlignFooter: alnFooter,
		AutoFormattingHeaders: autoFormattingHeaders,
		AutoMergeCells: autoMergeCells,
		AutoWrap: autoWrap,
		Caption: caption,
		CenterSeparator: centerSeparator,
		ColumnSeparator: columnSeparator,
		ColWidth: colWidth,
		FirstRowAsHeader: firstRowAsHeader,
		Footer: footer,
		Header: header,
		Newline: newline,
		Reflow: reflow,
		RowLine: rowLine,
		RowSeparator: rowSeparator,
	}

	fmt.Printf("%#v\n", renderOpts)
}

func Render(r io.Reader, w io.Writer, opts RenderOptions) error {
	csvR := csv.NewReader(r)
	tblW := tw.NewWriter(w)

	ConfigTableWriter(tblW, opts)

	ch := make(chan []string)
	done := make(chan struct{})

	go func() {
		for {
			rec, err := csvR.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				break
			}

			ch <- rec
		}

		close(ch)
	}()

	go func() {
		first := true
		for rec := range ch {
			if first {
				first = false
				tblW.SetHeader(rec)
				continue
			}
			tblW.Append(rec)
		}

		tblW.Render()
		close(done)
	}()

	<-done

	return nil
}

func ConfigTableWriter(t *tw.Table, opts RenderOptions) {
	if len(opts.AlignColumns) > 0 {
		var alignments []int
		for _, a := range opts.AlignColumns {
			alignments = append(alignments, a.ToTableWriter())
		}
		t.SetColumnAlignment(alignments)
	} else {
		t.SetAlignment(opts.Align.ToTableWriter())
	}

	if opts.AlignHeader != AlignDefault {
		t.SetHeaderAlignment(opts.AlignHeader.ToTableWriter())
	}

	if opts.AlignFooter != AlignDefault {
		t.SetFooterAlignment(opts.AlignFooter.ToTableWriter())
	}

	t.SetAutoFormatHeaders(opts.AutoFormattingHeaders)
	t.SetAutoMergeCells(opts.AutoMergeCells)
	t.SetAutoWrapText(opts.AutoWrap)

	if len(opts.Caption) > 0 {
		t.SetCaption(true, opts.Caption)
	}

	t.SetCenterSeparator(opts.CenterSeparator)
	t.SetColumnSeparator(opts.ColumnSeparator)

	if len(opts.Footer) > 0 {
		t.SetFooter(opts.Footer)
	}

	if len(opts.Header) > 0 {
		fmt.Println("testing")
		t.SetHeader(opts.Header)
	}

	t.SetNewLine(opts.Newline)
	t.SetReflowDuringAutoWrap(opts.Reflow)
	t.SetRowSeparator(opts.RowSeparator)
}
