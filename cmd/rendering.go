package cmd

import (
	"encoding/csv"
	tw "github.com/olekukonko/tablewriter"
	"io"
)

// RenderOptions defines the various rendering options available
type RenderOptions struct {
	Align                 Alignment
	AlignColumns          []Alignment
	AlignHeader           Alignment
	AlignFooter           Alignment
	AutoFormattingHeaders bool
	AutoMergeCells        bool
	AutoWrap              bool
	Caption               string
	CenterSeparator       string
	ColumnSeparator       string
	ColWidth              int
	FirstRowAsHeader      bool
	Footer                []string
	Header                []string
	Newline               string
	Reflow                bool
	RowLine               bool
	RowSeparator          string
}

// Render will read from r and write to w using opts to control how the table is rendered. Will return
// an error if r.Read() fails.
func Render(r *csv.Reader, w io.Writer, opts RenderOptions) error {
	tblW := tw.NewWriter(w)

	ConfigTableWriter(tblW, opts)

	// TODO: Test if a buffered channel is worth it.
	ch := make(chan []string)
	done := make(chan error)

	// Read from r and write to the channel
	go func() {
		defer close(ch)

		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				done <- err
				break
			}
			ch <- rec
		}
	}()

	// Read from the channel and write to the io.Writer
	go func() {
		defer close(done)

		first := true
		for rec := range ch {
			// if this is the first record and FirstRowAsHeader is set, use the
			// first row as the header.
			if first && opts.FirstRowAsHeader {
				first = false
				tblW.SetHeader(rec)
				continue
			}
			tblW.Append(rec)
		}
	}()

	err := <-done
	<-ch

	if err != nil {
		tblW.Render()
	}

	return err
}
