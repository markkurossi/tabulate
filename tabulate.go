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

// Align specifies cell alignment in horizontal and vertical
// directions.
type Align int

// Alignment constants. The first character specifies the vertical
// alignment (Top, Middle, Bottom) and the second character specifies
// the horizointal alignment (Left, Center, Right).
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

// Border specifies the table border drawing elements.
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

// WhiteSpace defines tabulation with whitespace (non-existing)
// borders.
var WhiteSpace = Border{}

// ASCII uses ASCII characters '-', '+', and '|' to draw the table
// borders.
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

// Unicode uses Unicode line drawing characters to draw the table
// borders.
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

// Colon uses the ':' character to mark vertical lines between cells.
var Colon = Border{
	VM: " : ",
}

// CSV defines the RFC 4180 Comma-Separated Values tabulation.
var CSV = Border{
	VM: ",",
	VR: "\r",
}

// Tabulate defined a tabulator instance.
type Tabulate struct {
	Padding int
	Border  Border
	Escape  Escape
	Headers []*Column
	Rows    []*Row
	asData  Data
}

// Escape is an escape function for converting table cell value into
// the output format.
type Escape func(string) string

// NewWS creates a new tabulator with the WhiteSpace borders.
func NewWS() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  WhiteSpace,
	}
}

// NewASCII creates a new tabulator with the ASCII borders.
func NewASCII() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  ASCII,
	}
}

// NewUnicode creates a new tabulator with the Unicode borders.
func NewUnicode() *Tabulate {
	return &Tabulate{
		Padding: 2,
		Border:  Unicode,
	}
}

// NewColon creates a new tabulator with the Colon borders.
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

// NewCSV creates a new tabulator for CVS outputs. It uses the CSV
// borders and an escape function which handles ',' and '\n'
// characters inside cell values.
func NewCSV() *Tabulate {
	return &Tabulate{
		Padding: 0,
		Border:  CSV,
		Escape:  escapeCSV,
	}
}

// Header adds a new column to the table and specifies its header
// label.
func (t *Tabulate) Header(label string) *Column {
	return t.HeaderData(NewLines(label))
}

// HeaderData adds a new column to the table and specifies is header
// data.
func (t *Tabulate) HeaderData(data Data) *Column {
	col := &Column{
		Data: data,
	}
	t.Headers = append(t.Headers, col)
	return col
}

// Row adds a new data row to the table.
func (t *Tabulate) Row() *Row {
	row := &Row{
		Tab: t,
	}
	t.Rows = append(t.Rows, row)
	return row
}

// Print layouts the table into the argument io.Writer.
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
			t.printColumn(o, hdr, idx, line, width, height)
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
				t.printColumn(o, col, idx, line, width, height)
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

func (t *Tabulate) printColumn(o io.Writer, col *Column,
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

func (t *Tabulate) data() Data {
	if t.asData == nil {
		builder := new(strings.Builder)
		t.Print(builder)
		t.asData = NewLines(builder.String())
	}
	return t.asData
}

// Width implements the Width of the Data interface.
func (t *Tabulate) Width() int {
	return t.data().Width()
}

// Height implements the Height of the Data interface.
func (t *Tabulate) Height() int {
	return t.data().Height()
}

// Content implements the Content of the Data interface.
func (t *Tabulate) Content(row int) string {
	return t.data().Content(row)
}

func (t *Tabulate) String() string {
	return t.data().String()
}

// Clone creates a new tabulator sharing the headers and their
// attributes. The new tabulator does not share the data rows with the
// original tabulator.
func (t *Tabulate) Clone() *Tabulate {
	return &Tabulate{
		Padding: t.Padding,
		Border:  t.Border,
		Escape:  t.Escape,
		Headers: t.Headers,
	}
}

// Row defines a data row in the tabulator.
type Row struct {
	Tab     *Tabulate
	Columns []*Column
}

// Height returns the row height in lines.
func (r *Row) Height() int {
	var max int
	for _, col := range r.Columns {
		if col.Data.Height() > max {
			max = col.Data.Height()
		}
	}
	return max
}

// Column adds a new string column to the row.
func (r *Row) Column(label string) *Column {
	return r.ColumnData(NewLines(label))
}

// ColumnData adds a new data column to the row.
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

// Column defines a table column data and its attributes.
type Column struct {
	Align  Align
	Data   Data
	Format Format
}

// SetAlign sets the column alignment.
func (col *Column) SetAlign(align Align) *Column {
	col.Align = align
	return col
}

// SetFormat sets the column text format.
func (col *Column) SetFormat(format Format) *Column {
	col.Format = format
	return col
}

// Width returns the column width in runes.
func (col *Column) Width() int {
	if col.Data == nil {
		return 0
	}
	return col.Data.Width()
}

// Height returns the column heigh in lines.
func (col *Column) Height() int {
	if col.Data == nil {
		return 0
	}
	return col.Data.Height()
}

// Content returns the specified row from the column. If the column
// does not have that many row, the function returns an empty string.
func (col *Column) Content(row int) string {
	if col.Data == nil {
		return ""
	}
	return col.Data.Content(row)
}

// Data contains table cell data.
type Data interface {
	Width() int
	Height() int
	Content(row int) string
	String() string
}

// Lines implements the Data interface over an array of lines.
type Lines struct {
	MaxWidth int
	Lines    []string
}

// NewLines creates a new Lines data from the argument string. The
// argument string is split into lines from the newline ('\n')
// character.
func NewLines(str string) *Lines {
	return NewLinesData(strings.Split(strings.TrimRight(str, "\n"), "\n"))
}

// NewLinesData creates a new Lines data from the array of strings.
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

// NewText creates a new Lines data, containing one line.
func NewText(str string) *Lines {
	return &Lines{
		MaxWidth: len([]rune(str)),
		Lines:    []string{str},
	}
}

// Width implements the Data.Width().
func (lines *Lines) Width() int {
	return lines.MaxWidth
}

// Height implements the Data.Height().
func (lines *Lines) Height() int {
	return len(lines.Lines)
}

// Content implements the Data.Content().
func (lines *Lines) Content(row int) string {
	if row >= lines.Height() {
		return ""
	}
	return lines.Lines[row]
}

func (lines *Lines) String() string {
	return strings.Join(lines.Lines, "\n")
}
