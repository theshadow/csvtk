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
	"io"
	"os"

	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"bufio"
	"encoding/csv"
)

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
    - tail -n +1 contact-list.csv | csv render --header "FName LName Email Phone"

  Using the first line in the CSV stream as an inline header.
    - cat contact-list.csv | csv render --first-row-as-header
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		align, err := cmd.Flags().GetString("align")
		if err != nil {
			return err
		}

		aln, err := FromString(align)
		if err != nil {
			return err
		}

		alignCols, err := cmd.Flags().GetStringSlice("align-columns")
		if err != nil {
			return err
		}

		var alnCols []Alignment
		err = FromStringArray(alignCols, alnCols)
		if err != nil {
			return err
		}

		alignHeader, err := cmd.Flags().GetString("align-header")
		if err != nil {
			return err
		}

		var alnHeader Alignment
		alnHeader, err = FromString(alignHeader)

		alignFooter, err := cmd.Flags().GetString("align-footer")
		if err != nil {
			return err
		}

		var alnFooter Alignment
		alnFooter, err = FromString(alignFooter)

		autoFormattingHeaders, err := cmd.Flags().GetBool("auto-formatting-headers")
		if err != nil {
			return err
		}

		autoMergeCells, err := cmd.Flags().GetBool("auto-merge-cells")
		if err != nil {
			return err
		}

		autoWrap, err := cmd.Flags().GetBool("auto-wrap")
		if err != nil {
			return err
		}

		caption, err := cmd.Flags().GetString("caption")
		if err != nil {
			return err
		}

		centerSeparator, err := cmd.Flags().GetString("center-separator")
		if err != nil {
			return err
		}

		columnSeparator, err := cmd.Flags().GetString("column-separator")
		if err != nil {
			return err
		}

		colWidth, err := cmd.Flags().GetInt("col-width")
		if err != nil {
			return err
		}

		firstRowAsHeader, err := cmd.Flags().GetBool("first-row-as-header")
		if err != nil {
			return err
		}

		footer, err := cmd.Flags().GetStringSlice("footer")
		if err != nil {
			return err
		}

		header, err := cmd.Flags().GetStringSlice("header")
		if err != nil {
			return err
		}

		newline, err := cmd.Flags().GetString("newline")
		if err != nil {
			return err
		}

		reflow, err := cmd.Flags().GetBool("reflow")
		if err != nil {
			return err
		}

		rowLine, err := cmd.Flags().GetBool("row-line")
		if err != nil {
			return err
		}

		rowSeparator, err := cmd.Flags().GetString("row-separator")
		if err != nil {
			return err
		}

		input, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		renderOpts := RenderOptions{
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

		var in io.Reader = os.Stdin
		var out io.WriteCloser = os.Stdout

		if len(input) > 0 {
			fIn, err := os.Open(input)
			defer fIn.Close()

			if err != nil {
				return err
			}

			in = bufio.NewReader(fIn)
		}

		if len(output) > 0 {
			fOut, err := os.Create(output)
			out = NewBufWriteCloser(fOut, bufio.NewWriter(fOut))
			if err != nil {
				return err
			}
		}

		csvR := csv.NewReader(in)

		err = Render(csvR, out, renderOpts)
		out.Close()
		return err
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringP("align", "a", string(AlignDefault),
		"Defines the text alignment for the entire table.")
	renderCmd.Flags().StringSliceP("align-columns", "", []string{},
	"Defines the alignment for each column individually.")
	renderCmd.Flags().String("align-header", string(AlignDefault),
		"Defines the alignment for the header.")
	renderCmd.Flags().String("align-footer", string(AlignDefault),
		"Defines the alignment for the header columns.")
	renderCmd.Flags().BoolP("auto-formatting-headers", "", true,
		"When specified auto formatting of the headers will be disabled.")
	renderCmd.Flags().BoolP("auto-merge-cells", "", false,
		"Defines the alignment for each footer column individually.")
	renderCmd.Flags().Bool("auto-wrap", true,
		"When set the text will not be automatically wrapped.")
	renderCmd.Flags().StringP("caption", "c", "",
		"Defines the caption message to be displayed with the table.")
	renderCmd.Flags().String("center-separator", tw.CENTER,
		"Defines what character will separate center columns.")
	renderCmd.Flags().String("column-separator", tw.COLUMN,
		"Defines what character will separate each column.")
	renderCmd.Flags().IntP("col-width", "w", tw.MAX_ROW_WIDTH,
		"Defines the fixed width for each column.")
	renderCmd.Flags().Bool("first-row-as-header", false,
		"When specified the first row will be treated as the headers.")
	renderCmd.Flags().StringSliceP("footer", "", []string{},
	"Defines what the footer columns should be.")
	renderCmd.Flags().StringSlice("header", []string{},
	"Defines what the header columns should be.")
	renderCmd.Flags().StringP("input", "i", "",
		"Define a file to read from instead of Stdin.")
	renderCmd.Flags().StringP("newline", "", tw.NEWLINE,
		"Defines what character will be used at the end of a line.")
	renderCmd.Flags().BoolP("reflow", "", true,
		"When specified the text will not be automatically re-flowed.")
	renderCmd.Flags().StringP("output", "o", "",
		"Define a file to write to instead of Stdout.")
	renderCmd.Flags().BoolP("row-line", "", false,
		"When specified each row will be delimited with a row line.")
	renderCmd.Flags().StringP("row-separator", "", tw.ROW,
		"Defines the character used to separate each row.")
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
		t.SetHeader(opts.Header)
	}

	t.SetNewLine(opts.Newline)
	t.SetReflowDuringAutoWrap(opts.Reflow)
	t.SetRowSeparator(opts.RowSeparator)
}

// Provides a WriteCloser for a bufio.Writer.
type BufWriteCloser struct {
	f *os.File
	*bufio.Writer
}

func NewBufWriteCloser(f *os.File, buf *bufio.Writer) io.WriteCloser {
	return &BufWriteCloser{f, buf}
}

func (b *BufWriteCloser) Close() error {
	if err := b.Flush(); err != nil {
		return err
	}
	return b.f.Close()
}