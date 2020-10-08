//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"encoding/json"
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

var aligns = map[Align]string{
	TL:   "TL",
	TC:   "TC",
	TR:   "TR",
	ML:   "ML",
	MC:   "MC",
	MR:   "MR",
	BL:   "BL",
	BC:   "BC",
	BR:   "BR",
	None: "None",
}

func (a Align) String() string {
	name, ok := aligns[a]
	if ok {
		return name
	}
	return fmt.Sprintf("{align %d}", a)
}

// Style specifies the table borders and rendering style.
type Style int

// Table styles.
const (
	Plain Style = iota
	ASCII
	Unicode
	UnicodeLight
	UnicodeBold
	Colon
	Simple
	Github
	CSV
	JSON
)

// Border specifies the table border drawing elements.
type Border struct {
	HT string
	HM string
	HB string
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

// Borders specifies the thable border drawing elements for the table
// header and body.
type Borders struct {
	Header Border
	Body   Border
}

var asciiBorder = Border{
	HT: "-",
	HM: "-",
	HB: "-",
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

var unicodeLight = Border{
	HT: "\u2500",
	HM: "\u2500",
	HB: "\u2500",
	VL: "\u2502",
	VM: "\u2502",
	VR: "\u2502",
	TL: "\u250C",
	TM: "\u252c",
	TR: "\u2510",
	ML: "\u251C",
	MM: "\u253C",
	MR: "\u2524",
	BL: "\u2514",
	BM: "\u2534",
	BR: "\u2518",
}

var unicodeBold = Border{
	HT: "\u2501",
	HM: "\u2501",
	HB: "\u2501",
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

var borders = map[Style]Borders{
	Plain: {},
	ASCII: {
		Header: asciiBorder,
		Body:   asciiBorder,
	},
	Unicode: {
		Header: Border{
			HT: "\u2501",
			HM: "\u2501",
			HB: "\u2501",
			VL: "\u2503",
			VM: "\u2503",
			VR: "\u2503",
			TL: "\u250F",
			TM: "\u2533",
			TR: "\u2513",
			ML: "\u2521",
			MM: "\u2547",
			MR: "\u2529",
			BL: "\u2517",
			BM: "\u253B",
			BR: "\u251B",
		},
		Body: Border{
			HT: "\u2500",
			HM: "\u2500",
			HB: "\u2500",
			VL: "\u2502",
			VM: "\u2502",
			VR: "\u2502",
			TL: "\u250C",
			TM: "\u252c",
			TR: "\u2510",
			ML: "\u251C",
			MM: "\u253C",
			MR: "\u2524",
			BL: "\u2514",
			BM: "\u2534",
			BR: "\u2518",
		},
	},
	UnicodeLight: {
		Header: unicodeLight,
		Body:   unicodeLight,
	},
	UnicodeBold: {
		Header: unicodeBold,
		Body:   unicodeBold,
	},
	Colon: {
		Header: Border{
			VM: " : ",
		},
		Body: Border{
			VM: " : ",
		},
	},
	Simple: {
		Header: Border{
			HM: "-",
			VM: "  ",
			MM: "  ",
		},
		Body: Border{
			VM: "  ",
			MM: "  ",
		},
	},
	Github: {
		Header: Border{
			HM: "-",
			VL: "|",
			VM: "|",
			VR: "|",
			ML: "|",
			MM: "|",
			MR: "|",
		},
		Body: Border{
			VL: "|",
			VM: "|",
			VR: "|",
		},
	},
	CSV: {
		Header: Border{
			VM: ",",
			VR: "\r",
		},
		Body: Border{
			VM: ",",
			VR: "\r",
		},
	},
	JSON: {},
}

// Tabulate defined a tabulator instance.
type Tabulate struct {
	Padding int
	Borders Borders
	Escape  Escape
	Output  func(t *Tabulate, o io.Writer)
	Headers []*Column
	Rows    []*Row
	asData  Data
}

// Escape is an escape function for converting table cell value into
// the output format.
type Escape func(string) string

// New creates a new tabulate object with the specified rendering
// style.
func New(style Style) *Tabulate {
	tab := &Tabulate{
		Padding: 2,
		Borders: borders[style],
	}
	switch style {
	case Colon, Simple:
		tab.Padding = 0
	case CSV:
		tab.Padding = 0
		tab.Escape = escapeCSV
	case JSON:
		tab.Padding = 0
		tab.Output = outputJSON
	}
	return tab
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

func outputJSON(t *Tabulate, o io.Writer) {
	data, err := json.Marshal(t)
	if err != nil {
		fmt.Fprintf(o, "JSON marshal failed: %s", err)
		return
	}
	fmt.Fprintf(o, string(data))
	fmt.Fprintln(o)
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
	if t.Output != nil {
		t.Output(t, o)
		return
	}
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
	if len(t.Borders.Header.HT) > 0 {
		fmt.Fprint(o, t.Borders.Header.TL)
		for idx, width := range widths {
			for i := 0; i < width+t.Padding; i++ {
				fmt.Fprint(o, t.Borders.Header.HT)
			}
			if idx+1 < len(widths) {
				fmt.Fprint(o, t.Borders.Header.TM)
			} else {
				fmt.Fprintln(o, t.Borders.Header.TR)
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
			t.printColumn(o, true, hdr, idx, line, width, height)
		}
		fmt.Fprintln(o, t.Borders.Header.VR)
	}

	if len(t.Borders.Header.HM) > 0 {
		fmt.Fprint(o, t.Borders.Header.ML)
		for idx, width := range widths {
			for i := 0; i < width+t.Padding; i++ {
				fmt.Fprint(o, t.Borders.Header.HM)
			}
			if idx+1 < len(widths) {
				fmt.Fprint(o, t.Borders.Header.MM)
			} else {
				fmt.Fprintln(o, t.Borders.Header.MR)
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
				t.printColumn(o, false, col, idx, line, width, height)
			}
			fmt.Fprintln(o, t.Borders.Body.VR)
		}
	}

	if len(t.Borders.Body.HB) > 0 {
		fmt.Fprint(o, t.Borders.Body.BL)
		for idx, width := range widths {
			for i := 0; i < width+t.Padding; i++ {
				fmt.Fprint(o, t.Borders.Body.HB)
			}
			if idx+1 < len(widths) {
				fmt.Fprint(o, t.Borders.Body.BM)
			} else {
				fmt.Fprintln(o, t.Borders.Body.BR)
			}
		}
	}
}

func (t *Tabulate) printColumn(o io.Writer, hdr bool, col *Column,
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

	if hdr {
		if idx == 0 {
			fmt.Fprint(o, t.Borders.Header.VL)
		} else {
			fmt.Fprint(o, t.Borders.Header.VM)
		}
	} else {
		if idx == 0 {
			fmt.Fprint(o, t.Borders.Body.VL)
		} else {
			fmt.Fprint(o, t.Borders.Body.VM)
		}
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

// Width implements the Data.Width().
func (t *Tabulate) Width() int {
	return t.data().Width()
}

// Height implements the Data.Height().
func (t *Tabulate) Height() int {
	return t.data().Height()
}

// Content implements the Data.Content().
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
		Borders: t.Borders,
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
