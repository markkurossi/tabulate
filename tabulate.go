//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"golang.org/x/text/width"
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
	CompactUnicode
	CompactUnicodeLight
	CompactUnicodeBold
	Colon
	Simple
	SimpleUnicode
	SimpleUnicodeBold
	Github
	CSV
	JSON
)

// Styles list all supported tabulation types.
var Styles = map[string]Style{
	"plain":          Plain,
	"ascii":          ASCII,
	"uc":             Unicode,
	"uclight":        UnicodeLight,
	"ucbold":         UnicodeBold,
	"compactuc":      CompactUnicode,
	"compactuclight": CompactUnicodeLight,
	"compactucbold":  CompactUnicodeBold,
	"colon":          Colon,
	"simple":         Simple,
	"simpleuc":       SimpleUnicode,
	"simpleucbold":   SimpleUnicodeBold,
	"github":         Github,
	"csv":            CSV,
	"json":           JSON,
}

func (s Style) String() string {
	for name, style := range Styles {
		if s == style {
			return name
		}
	}
	return fmt.Sprintf("{Style %d}", s)
}

// StyleNames returns the tabulation style names as a sorted slice.
func StyleNames() []string {
	var names []string
	for name := range Styles {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })
	return names
}

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

var unicodeHeader = Border{
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
}

var unicodeBody = Border{
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
		Header: unicodeHeader,
		Body:   unicodeBody,
	},
	UnicodeLight: {
		Header: unicodeLight,
		Body:   unicodeLight,
	},
	UnicodeBold: {
		Header: unicodeBold,
		Body:   unicodeBold,
	},
	CompactUnicode: {
		Header: unicodeHeader,
		Body:   unicodeBody,
	},
	CompactUnicodeLight: {
		Header: unicodeLight,
		Body:   unicodeLight,
	},
	CompactUnicodeBold: {
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
			VM: " ",
			MM: " ",
		},
		Body: Border{
			VM: " ",
			MM: " ",
		},
	},
	SimpleUnicode: {
		Header: Border{
			HM: "\u2500",
			VM: " ",
			MM: " ",
		},
		Body: Border{
			VM: " ",
			MM: " ",
		},
	},
	SimpleUnicodeBold: {
		Header: Border{
			HM: "\u2501",
			VM: " ",
			MM: " ",
		},
		Body: Border{
			VM: " ",
			MM: " ",
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
	Padding     int
	TrimColumns bool
	Borders     Borders
	Measure     Measure
	Escape      Escape
	Output      func(t *Tabulate, o io.Writer)
	Defaults    []Align
	Headers     []*Column
	Rows        []*Row
	asData      Data
}

// Measure returns the column width in display units. This can be used
// to remove any non-printable formatting codes from the value.
type Measure func(column string) int

// MeasureRunes measures the column width by counting its runes. This
// assumes that all runes have the same width consuming single output
// column cell.
func MeasureRunes(column string) int {
	return len([]rune(column))
}

// MeasureUnicode measures the column width by taking into
// consideration East Asian Wide characters. The function assumes that
// East Asian Wide characters consume two output column cells.
func MeasureUnicode(column string) int {
	var w int
	for _, r := range column {
		if width.LookupRune(r).Kind() == width.EastAsianWide {
			w += 2
		} else {
			w++
		}
	}
	return w
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
		Measure: MeasureUnicode,
	}
	switch style {
	case Colon, Simple, SimpleUnicode, SimpleUnicodeBold,
		CompactUnicode, CompactUnicodeLight, CompactUnicodeBold:
		tab.Padding = 0
	case CSV:
		tab.Padding = 0
		tab.TrimColumns = true
		tab.Escape = escapeCSV
	case JSON:
		tab.Padding = 0
		tab.TrimColumns = true
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

// SetDefaults sets the column default attributes. These are used if
// the table does not have headers.
func (t *Tabulate) SetDefaults(col int, align Align) {
	for len(t.Defaults) <= col {
		t.Defaults = append(t.Defaults, TL)
	}
	t.Defaults[col] = align
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
	if len(t.Headers) == 0 && len(t.Rows) == 0 {
		// No columns to tabulate.
		return
	}
	if t.Output != nil {
		t.Output(t, o)
		return
	}
	// Measure columns.
	widths := make([]int, len(t.Headers))
	for idx, hdr := range t.Headers {
		w := hdr.Data.Width(t.Measure)
		if w > widths[idx] {
			widths[idx] = w
		}
	}
	for _, row := range t.Rows {
		for idx, col := range row.Columns {
			if idx >= len(widths) {
				widths = append(widths, 0)
			}
			w := col.Width(t.Measure)
			if w > widths[idx] {
				widths[idx] = w
			}
		}
	}

	if len(t.Headers) > 0 {
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
	}

	var bottomBorder Border

	if len(t.Rows) > 0 {
		if len(t.Headers) > 0 {
			// Both headers and rows.
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
		} else {
			// Only rows.
			if len(t.Borders.Body.HT) > 0 {
				fmt.Fprint(o, t.Borders.Body.TL)
				for idx, width := range widths {
					for i := 0; i < width+t.Padding; i++ {
						fmt.Fprint(o, t.Borders.Body.HT)
					}
					if idx+1 < len(widths) {
						fmt.Fprint(o, t.Borders.Body.TM)
					} else {
						fmt.Fprintln(o, t.Borders.Body.TR)
					}
				}
			}
		}

		// Data rows.
		for _, row := range t.Rows {
			height := row.Height()

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
		// Use the body graphics to close the table.
		bottomBorder = t.Borders.Body
	} else {
		// No data rows. Use the header graphics to close the table.
		bottomBorder = t.Borders.Header
	}

	if len(bottomBorder.HB) > 0 {
		fmt.Fprint(o, bottomBorder.BL)
		for idx, width := range widths {
			for i := 0; i < width+t.Padding; i++ {
				fmt.Fprint(o, bottomBorder.HB)
			}
			if idx+1 < len(widths) {
				fmt.Fprint(o, bottomBorder.BM)
			} else {
				fmt.Fprintln(o, bottomBorder.BR)
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

	pad := width - t.Measure(content)
	if t.TrimColumns {
		pad = 0
	}
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
func (t *Tabulate) Width(m Measure) int {
	return t.data().Width(m)
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
		Padding:     t.Padding,
		TrimColumns: t.TrimColumns,
		Borders:     t.Borders,
		Measure:     t.Measure,
		Escape:      t.Escape,
		Defaults:    t.Defaults,
		Headers:     t.Headers,
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
	} else if idx < len(r.Tab.Defaults) {
		hdr = &Column{
			Align: r.Tab.Defaults[idx],
		}
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
func (col *Column) Width(m Measure) int {
	if col.Data == nil {
		return 0
	}
	return col.Data.Width(m)
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
