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

type Align int

const (
	TL Align = iota
	TC
	TR
	ML
	MC
	MR
	BL
	BC
	BR
	None
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
	Headers []*Column
	Rows    []*Row
}

type Escape func(string) string

func NewWS() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  WhiteSpace,
	}
}

func NewASCII() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  ASCII,
	}
}

func NewUnicode() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  Unicode,
	}
}

func NewColon() *Tabulate {
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

func NewCSV() *Tabulate {
	return &Tabulate{
		Padding: 0,
		Border:  CSV,
		Escape:  escapeCSV,
	}
}

func (t *Tabulate) Header(label string) *Column {
	return t.HeaderData(NewLines(label))
}

func (t *Tabulate) HeaderData(data Data) *Column {
	col := &Column{
		Data: data,
	}
	t.Headers = append(t.Headers, col)
	return col
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
			var hdr *Column
			if idx < len(t.Headers) {
				hdr = t.Headers[idx]
			} else {
				hdr = &Column{}
			}
			t.PrintColumn(o, hdr, idx, line, width, height)
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
				var col *Column
				if idx < len(row.Columns) {
					col = row.Columns[idx]
				} else {
					col = &Column{}
				}
				t.PrintColumn(o, col, idx, line, width, height)
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

func (t *Tabulate) PrintColumn(o io.Writer, col *Column,
	idx, line, width, height int) {

	vspace := height - col.Height()
	switch col.Align {
	case TL, TC, TR, None:

	case ML, MC, MR:
		line -= vspace / 2

	case BL, BC, BR:
		line -= vspace
	}

	var content string
	if line >= 0 {
		content = col.Content(line)
	}
	if t.Escape != nil {
		content = t.Escape(content)
	}

	lPad := t.Padding / 2
	rPad := t.Padding - lPad

	pad := width - len([]rune(content))
	switch col.Align {
	case None:
		lPad = 0
		rPad = 0

	case TL, ML, BL:
		rPad += pad

	case TC, MC, BC:
		l := pad / 2
		r := pad - l
		lPad += l
		rPad += r

	case TR, MR, BR:
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

func (t *Tabulate) Clone() *Tabulate {
	return &Tabulate{
		Padding: t.Padding,
		Border:  t.Border,
		Escape:  t.Escape,
		Headers: t.Headers,
	}
}

type Row struct {
	Tab     *Tabulate
	Columns []*Column
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

func (r *Row) Column(label string) *Column {
	return r.ColumnData(NewLines(label))
}

func (r *Row) ColumnData(data Data) *Column {
	idx := len(r.Columns)
	var hdr *Column
	if idx < len(r.Tab.Headers) {
		hdr = r.Tab.Headers[idx]
	} else {
		hdr = &Column{}
	}

	col := &Column{
		Align:  hdr.Align,
		Data:   data,
		Format: hdr.Format,
	}

	r.Columns = append(r.Columns, col)
	return col
}

type Column struct {
	Align  Align
	Data   Data
	Format Format
}

func (col *Column) SetAlign(align Align) *Column {
	col.Align = align
	return col
}

func (col *Column) SetFormat(format Format) *Column {
	col.Format = format
	return col
}

func (col *Column) Width() int {
	if col.Data == nil {
		return 0
	}
	return col.Data.Width()
}

func (col *Column) Height() int {
	if col.Data == nil {
		return 0
	}
	return col.Data.Height()
}

func (col *Column) Content(row int) string {
	if col.Data == nil {
		return ""
	}
	return col.Data.Content(row)
}

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
	return NewLinesData(strings.Split(strings.TrimRight(str, "\n"), "\n"))
}

func NewLinesData(lines []string) *Lines {
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
