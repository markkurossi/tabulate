//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"fmt"
	"io"
	"strings"
)

type Border struct {
	H  string
	V  string
	TL string
	TM string
	TR string
	ML string
	MM string
	MR string
	BL string
	BM string
	BR string
}

var WhiteSpace = Border{}

var ASCII = Border{
	H:  "-",
	V:  "|",
	TL: "+",
	TM: "+",
	TR: "+",
	ML: "+",
	MM: "+",
	MR: "+",
	BL: "+",
	BM: "+",
	BR: "+",
}

var Unicode = Border{
	H:  "\u2501",
	V:  "\u2503",
	TL: "\u250F",
	TM: "\u2533",
	TR: "\u2513",
	ML: "\u2523",
	MM: "\u254B",
	MR: "\u252B",
	BL: "\u2517",
	BM: "\u253B",
	BR: "\u251B",
}

type Tabulate struct {
	Padding int
	Border  Border
	Headers []Column
	Rows    []*Row
}

func NewTabulateWS() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  WhiteSpace,
	}
}

func NewTabulateASCII() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  ASCII,
	}
}

func NewTabulateUnicode() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  Unicode,
	}
}

func (t *Tabulate) Header(align Align, data Data) *Tabulate {
	t.Headers = append(t.Headers, Column{
		Align: align,
		Data:  data,
	})
	return t
}

func (t *Tabulate) Row() *Row {
	row := &Row{
		Tab: t,
	}
	t.Rows = append(t.Rows, row)
	return row
}

func (t *Tabulate) Print(o io.Writer) {
	widths := make([]int, len(t.Headers))

	for idx, hdr := range t.Headers {
		if hdr.Data.Width() > widths[idx] {
			widths[idx] = hdr.Data.Width()
		}
	}
	for _, row := range t.Rows {
		for idx, col := range row.Columns {
			if idx >= len(widths) {
				widths = append(widths, 0)
			}
			if col.Width() > widths[idx] {
				widths[idx] = col.Width()
			}
		}
	}

	// Header.
	if len(t.Border.H) > 0 {
		fmt.Fprint(o, t.Border.TL)
		for idx, width := range widths {
			for i := 0; i < width+t.Padding; i++ {
				fmt.Fprint(o, t.Border.H)
			}
			if idx+1 < len(widths) {
				fmt.Fprint(o, t.Border.TM)
			} else {
				fmt.Fprintln(o, t.Border.TR)
			}
		}
	}

	var height int
	for _, hdr := range t.Headers {
		if hdr.Data.Height() > height {
			height = hdr.Data.Height()
		}
	}
	for line := 0; line < height; line++ {
		for idx, width := range widths {
			var hdr Column
			if idx < len(t.Headers) {
				hdr = t.Headers[idx]
			}
			t.PrintColumn(o, hdr, line, width)
		}
		fmt.Fprintln(o, t.Border.V)
	}

	if len(t.Border.H) > 0 {
		fmt.Fprint(o, t.Border.ML)
		for idx, width := range widths {
			for i := 0; i < width+t.Padding; i++ {
				fmt.Fprint(o, t.Border.H)
			}
			if idx+1 < len(widths) {
				fmt.Fprint(o, t.Border.MM)
			} else {
				fmt.Fprintln(o, t.Border.MR)
			}
		}
	}

	// Data rows.
	for _, row := range t.Rows {
		height = row.Height()

		for line := 0; line < height; line++ {
			for idx, width := range widths {
				var col Column
				if idx < len(row.Columns) {
					col = row.Columns[idx]
				}
				t.PrintColumn(o, col, line, width)
			}
			fmt.Fprintln(o, t.Border.V)
		}
	}

	if len(t.Border.H) > 0 {
		fmt.Fprint(o, t.Border.BL)
		for idx, width := range widths {
			for i := 0; i < width+t.Padding; i++ {
				fmt.Fprint(o, t.Border.H)
			}
			if idx+1 < len(widths) {
				fmt.Fprint(o, t.Border.BM)
			} else {
				fmt.Fprintln(o, t.Border.BR)
			}
		}
	}
}

func (t *Tabulate) PrintColumn(o io.Writer, col Column, line, width int) {

	lPad := t.Padding / 2
	rPad := t.Padding - lPad

	content := col.Content(line)

	pad := width - len([]rune(content))
	switch col.Align {
	case AlignLeft:
		rPad += pad

	case AlignCenter:
		l := pad / 2
		r := pad - l
		lPad += l
		rPad += r

	case AlignRight:
		lPad += pad
	}

	fmt.Fprint(o, t.Border.V)
	for i := 0; i < lPad; i++ {
		fmt.Fprint(o, " ")
	}
	if col.Format != FmtNone {
		fmt.Fprint(o, col.Format.VT100())
	}
	fmt.Fprint(o, content)
	if col.Format != FmtNone {
		fmt.Fprint(o, FmtNone.VT100())
	}
	for i := 0; i < rPad; i++ {
		fmt.Fprint(o, " ")
	}
}

func (t *Tabulate) Data() Data {
	builder := new(strings.Builder)
	t.Print(builder)
	return NewLines(builder.String())
}

type Row struct {
	Tab     *Tabulate
	Columns []Column
}

func (r *Row) Height() int {
	var max int
	for _, col := range r.Columns {
		if col.Data.Height() > max {
			max = col.Data.Height()
		}
	}
	return max
}

func (r *Row) Column(data Data) {
	idx := len(r.Columns)
	var hdr Column
	if idx < len(r.Tab.Headers) {
		hdr = r.Tab.Headers[idx]
	}

	r.Columns = append(r.Columns, Column{
		Align:  hdr.Align,
		Data:   data,
		Format: hdr.Format,
	})
}

func (r *Row) ColumnAttrs(align Align, data Data, format Format) {
	r.Columns = append(r.Columns, Column{
		Align:  align,
		Data:   data,
		Format: format,
	})
}

type Column struct {
	Align  Align
	Data   Data
	Format Format
}

func (col Column) Width() int {
	if col.Data == nil {
		return 0
	}
	return col.Data.Width()
}

func (col Column) Content(row int) string {
	if col.Data == nil {
		return ""
	}
	return col.Data.Content(row)
}

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

type Data interface {
	Width() int
	Height() int
	Content(row int) string
}

type Lines struct {
	MaxWidth int
	Lines    []string
}

func NewLines(str string) *Lines {
	lines := strings.Split(strings.TrimSpace(str), "\n")

	var max int
	for _, line := range lines {
		l := len([]rune(line))
		if l > max {
			max = l
		}
	}

	return &Lines{
		MaxWidth: max,
		Lines:    lines,
	}
}

func (lines *Lines) Width() int {
	return lines.MaxWidth
}

func (lines *Lines) Height() int {
	return len(lines.Lines)
}

func (lines *Lines) Content(row int) string {
	if row >= lines.Height() {
		return ""
	}
	return lines.Lines[row]
}
