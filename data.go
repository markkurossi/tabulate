//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"fmt"
	"strings"
)

var (
	_ = Data((&Value{}))
	_ = Data((&Lines{}))
	_ = Data((&Slice{}))
)

// Data contains table cell data.
type Data interface {
	Width(m Measure) int
	Height() int
	Content(row int) string
	String() string
}

// Value implements the Data interface for single value, such as bool,
// integer, etc.
type Value struct {
	string string
	value  interface{}
}

// NewValue creates a new Value for the argument value element.
func NewValue(v interface{}) *Value {
	return &Value{
		string: fmt.Sprintf("%v", v),
		value:  v,
	}
}

// Width implements the Data.Width().
func (v *Value) Width(m Measure) int {
	return m(v.string)
}

// Height implements the Data.Height().
func (v *Value) Height() int {
	return 1
}

// Content implements the Data.Content().
func (v *Value) Content(row int) string {
	if row > 0 {
		return ""
	}
	return v.string
}

func (v *Value) String() string {
	return v.string
}

// Lines implements the Data interface over an array of lines.
type Lines struct {
	Lines []string
}

// NewLines creates a new Lines data from the argument string. The
// argument string is split into lines from the newline ('\n')
// character.
func NewLines(str string) *Lines {
	return NewLinesData(strings.Split(strings.TrimRight(str, "\n"), "\n"))
}

// NewLinesData creates a new Lines data from the array of strings.
func NewLinesData(lines []string) *Lines {
	return &Lines{
		Lines: lines,
	}
}

// NewText creates a new Lines data, containing one line.
func NewText(str string) *Lines {
	return &Lines{
		Lines: []string{str},
	}
}

// Width implements the Data.Width().
func (lines *Lines) Width(m Measure) int {
	var max int
	for _, l := range lines.Lines {
		w := m(l)
		if w > max {
			max = w
		}
	}
	return max
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

// NewSlice creates a new Slice Data type with the specified maximum
// rendering width.
func NewSlice(maxWidth int) *Slice {
	return &Slice{
		maxWidth: maxWidth,
	}
}

// Slice implements the Data interface for an array of Data elements.
type Slice struct {
	maxWidth int
	height   int
	content  []Data
	lines    []string
}

func (arr *Slice) addLine(line string) {
	arr.lines = append(arr.lines, line)
}

func (arr *Slice) layout() {
	if len(arr.lines) > 0 {
		return
	}
	var line string
	for _, c := range arr.content {
		h := c.Height()
		if h == 0 {
			continue
		}
		if h > 1 {
			if len(line) > 0 {
				arr.addLine(line)
				line = ""
			}
			for row := 0; row < h; row++ {
				arr.addLine(c.Content(row))
			}
		} else {
			l := c.Content(0)
			if len(line) == 0 {
				line = l
			} else if len(line)+len(l) <= arr.maxWidth {
				line += " "
				line += l
			} else {
				arr.addLine(line)
				line = l
			}
		}
	}
	if len(line) > 0 {
		arr.addLine(line)
	}
}

// Append adds data to the array.
func (arr *Slice) Append(data Data) {
	arr.content = append(arr.content, data)
}

// Width implements the Data.Width().
func (arr *Slice) Width(m Measure) int {
	arr.layout()

	var max int
	for _, l := range arr.lines {
		w := m(l)
		if w > max {
			max = w
		}
	}
	return max
}

// Height implements the Data.Height().
func (arr *Slice) Height() int {
	arr.layout()
	return len(arr.lines)
}

// Content implements the Data.Content().
func (arr *Slice) Content(row int) string {
	if row < len(arr.lines) {
		return arr.lines[row]
	}
	return ""
}

func (arr *Slice) String() string {
	result := "["
	for idx, c := range arr.content {
		if idx > 0 {
			result += ","
		}
		result += c.String()
	}
	return result + "]"
}
