//
// Copyright (c) 2020 Markku Rossi
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
	_ = Data((&Array{}))
)

// Data contains table cell data.
type Data interface {
	Width() int
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
func (v *Value) Width() int {
	return len([]rune(v.string))
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

// Array implements the Data interface for an array of Data elements.
type Array struct {
	maxWidth int
	height   int
	content  []Data
}

// Append adds data to the array.
func (arr *Array) Append(data Data) {
	w := data.Width()
	if w > arr.maxWidth {
		arr.maxWidth = w
	}
	arr.height += data.Height()
	arr.content = append(arr.content, data)
}

// Width implements the Data.Width().
func (arr *Array) Width() int {
	return arr.maxWidth
}

// Height implements the Data.Height().
func (arr *Array) Height() int {
	return arr.height
}

// Content implements the Data.Content().
func (arr *Array) Content(row int) string {
	for _, c := range arr.content {
		h := c.Height()
		if h > row {
			return c.Content(row)
		}
		row -= h
	}
	return ""
}

func (arr *Array) String() string {
	result := "["
	for idx, c := range arr.content {
		if idx > 0 {
			result += ","
		}
		result += c.String()
	}
	return result + "]"
}
