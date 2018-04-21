package cmd

import (
	"io"
	"encoding/csv"
	tw "github.com/olekukonko/tablewriter"
)

// Specifies the various render options
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

func Render(r *csv.Reader, w io.Writer, opts RenderOptions) error {
	tblW := tw.NewWriter(w)

	ConfigTableWriter(tblW, opts)

	ch := make(chan []string)
	done := make(chan struct{})

	go func() {
		defer close(ch)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			ch <- rec
		}
	}()

	go func() {
		defer close(done)

		first := true
		for rec := range ch {
			if first && opts.FirstRowAsHeader {
				first = false
				tblW.SetHeader(rec)
				continue
			}
			tblW.Append(rec)
		}
	}()

	<-done
	tblW.Render()

	return nil
}
