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
	VL string
	VM string
	VR string
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
	VL: "|",
	VM: "|",
	VR: "|",
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
	VL: "\u2503",
	VM: "\u2503",
	VR: "\u2503",
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

var Colon = Border{
	VM: " : ",
}

var CSV = Border{
	VM: ",",
	VR: "\r",
}

type Tabulate struct {
	Padding int
	Border  Border
	Escape  Escape
	Headers []Column
	Rows    []*Row
}

type Escape func(string) string

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

func NewTabulateColon() *Tabulate {
	return &Tabulate{
		Padding: 0,
		Border:  Colon,
	}
}

func escapeCSV(val string) string {
	idxQuote := strings.IndexRune(val, '"')
	idxNewline := strings.IndexRune(val, '\n')

	if idxQuote < 0 && idxNewline < 0 {
		return val
	}

	var runes []rune
	runes = append(runes, '"')
	for _, r := range []rune(val) {
		if r == '"' {
			runes = append(runes, '"')
		}
		runes = append(runes, r)
	}
	runes = append(runes, '"')

	return string(runes)
}

func NewTabulateCSV() *Tabulate {
	return &Tabulate{
		Padding: 0,
		Border:  CSV,
		Escape:  escapeCSV,
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
			t.PrintColumn(o, hdr, idx, line, width)
		}
		fmt.Fprintln(o, t.Border.VR)
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
				t.PrintColumn(o, col, idx, line, width)
			}
			fmt.Fprintln(o, t.Border.VR)
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

func (t *Tabulate) PrintColumn(o io.Writer, col Column, idx, line, width int) {

	lPad := t.Padding / 2
	rPad := t.Padding - lPad

	content := col.Content(line)
	if t.Escape != nil {
		content = t.Escape(content)
	}

	pad := width - len([]rune(content))
	switch col.Align {
	case AlignNone:
		lPad = 0
		rPad = 0

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

	if idx == 0 {
		fmt.Fprint(o, t.Border.VL)
	} else {
		fmt.Fprint(o, t.Border.VM)
	}
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
	AlignNone Align = iota
	AlignLeft
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

func NewText(str string) *Lines {
	return &Lines{
		MaxWidth: len([]rune(str)),
		Lines:    []string{str},
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
