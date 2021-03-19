//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"encoding"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Flags control how reflection tabulation operates on different
// values.
type Flags int

// Flag values for reflection tabulation.
const (
	OmitEmpty Flags = 1 << iota
	InheritHeaders
)

const nilLabel = "<nil>"

// Reflect tabulates the value into the tabulation object. The flags
// control how different values are handled. The tags lists element
// tags which are included in reflection. If the element does not have
// tabulation tag, then it is always included in tabulation.
func Reflect(tab *Tabulate, flags Flags, tags []string, v interface{}) error {
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag] = true
	}

	value := reflect.ValueOf(v)

	// Follows pointers.
	for value.Type().Kind() == reflect.Ptr {
		if value.IsZero() {
			return nil
		}
		value = reflect.Indirect(value)
	}

	if value.Type().Kind() == reflect.Struct {
		return reflectStruct(tab, flags, tagMap, value)
	}
	if value.Type().Kind() == reflect.Map {
		return reflectMap(tab, flags, tagMap, value)
	}

	data, err := reflectValue(tab, flags, tagMap, value)
	if err != nil {
		return err
	}
	row := tab.Row()
	row.Column("")
	row.ColumnData(data)

	return nil
}

// Array tabulates the argument v into rows and columns. If the tab
// defines header columns, those will be used. Otherwise the first row
// of v defines the header columns.
func Array(tab *Tabulate, v [][]interface{}) (*Tabulate, error) {
	flags := OmitEmpty
	tags := make(map[string]bool)

	if len(tab.Headers) == 0 {
		if len(v) == 0 {
			return tab, nil
		}
		for _, c := range v[0] {
			data, err := reflectValue(tab, flags, tags, reflect.ValueOf(c))
			if err != nil {
				return nil, err
			}
			tab.HeaderData(data)
		}
		v = v[1:]
	}
	for _, r := range v {
		row := tab.Row()
		for _, c := range r {
			data, err := reflectValue(tab, flags, tags, reflect.ValueOf(c))
			if err != nil {
				return nil, err
			}
			row.ColumnData(data)
		}
	}
	return tab, nil
}

func reflectValue(tab *Tabulate, flags Flags, tags map[string]bool,
	value reflect.Value) (Data, error) {

	if value.CanInterface() {
		switch v := value.Interface().(type) {
		case encoding.TextMarshaler:
			data, err := v.MarshalText()
			if err != nil {
				return nil, err
			}
			return NewLinesData([]string{string(data)}), nil
		}
	}

	// Resolve interfaces.
	for value.Type().Kind() == reflect.Interface {
		if value.IsZero() {
			if flags&OmitEmpty == 0 {
				return NewLinesData([]string{nilLabel}), nil
			}
			return NewLinesData(nil), nil
		}
		value = value.Elem()
	}

	// Follow pointers.
	for value.Type().Kind() == reflect.Ptr {
		if value.IsZero() {
			if flags&OmitEmpty == 0 {
				return NewLinesData([]string{nilLabel}), nil
			}
		}
		value = reflect.Indirect(value)
	}

	switch value.Type().Kind() {
	case reflect.Bool:
		return NewValue(value.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewValue(value.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return NewValue(value.Uint()), nil

	case reflect.Float32, reflect.Float64:
		return NewValue(value.Float()), nil

	case reflect.Map:
		if value.Len() > 0 || flags&OmitEmpty == 0 {
			sub := tab.Clone()
			if flags&InheritHeaders == 0 {
				sub.Headers = nil
			}
			err := reflectMap(sub, flags, tags, value)
			if err != nil {
				return nil, err
			}
			return sub, nil
		}
		return NewLinesData(nil), nil

	case reflect.String:
		text := value.String()
		lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
		return NewLinesData(lines), nil

	case reflect.Slice:
		// Check slice element type.
		switch value.Type().Elem().Kind() {
		case reflect.Uint8:
			return reflectByteSliceValue(tab, flags, tags, value)

		case reflect.Int, reflect.Uint:
			return reflectSliceValue(tab, flags, tags, 40, value)

		default:
			return reflectSliceValue(tab, flags, tags, 0, value)
		}

	case reflect.Struct:
		sub := tab.Clone()
		if flags&InheritHeaders == 0 {
			sub.Headers = nil
		}
		err := reflectStruct(sub, flags, tags, value)
		if err != nil {
			return nil, err
		}
		return sub, nil

	default:
		text := value.String()
		if len(text) == 0 && flags&OmitEmpty == 1 {
			return NewLinesData(nil), nil
		}
		return NewLinesData([]string{text}), nil
	}
}

func reflectByteSliceValue(tab *Tabulate, flags Flags, tags map[string]bool,
	value reflect.Value) (Data, error) {

	arr, ok := value.Interface().([]byte)
	if !ok {
		return nil, fmt.Errorf("reflectByteSliceValue called for %T",
			value.Type().Kind())
	}

	const lineLength = 32
	var lines []string
	for i := 0; i < len(arr); i += lineLength {
		l := len(arr) - i
		if l > lineLength {
			l = lineLength
		}
		lines = append(lines, fmt.Sprintf("%x", arr[i:i+l]))
	}
	return NewLinesData(lines), nil
}

func reflectSliceValue(tab *Tabulate, flags Flags, tags map[string]bool,
	width int, value reflect.Value) (Data, error) {

	data := NewSlice(width)
loop:
	for i := 0; i < value.Len(); i++ {
		v := value.Index(i)
		// Follow pointers.
		for v.Type().Kind() == reflect.Ptr {
			if v.IsZero() {
				if flags&OmitEmpty == 0 {
					data.Append(NewText(nilLabel))
				}
				continue loop
			}
			v = reflect.Indirect(v)
		}
		switch v.Type().Kind() {
		case reflect.Struct:
			sub := tab.Clone()
			if flags&InheritHeaders == 0 {
				sub.Headers = nil
			}
			err := reflectStruct(sub, flags, tags, v)
			if err != nil {
				return nil, err
			}
			data.Append(sub)

		default:
			sub, err := reflectValue(tab, flags, tags, v)
			if err != nil {
				return nil, err
			}
			data.Append(sub)
		}
	}

	return data, nil
}

type row struct {
	key Data
	val Data
}

func reflectMap(tab *Tabulate, flags Flags, tags map[string]bool,
	v reflect.Value) error {

	var rows []row
	iter := v.MapRange()
	for iter.Next() {
		keyData, err := reflectValue(tab, flags, tags, iter.Key())
		if err != nil {
			return err
		}
		valData, err := reflectValue(tab, flags, tags, iter.Value())
		if err != nil {
			return err
		}
		rows = append(rows, row{
			key: keyData,
			val: valData,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		di := rows[i].key
		dj := rows[j].key

		height := di.Height()
		if dj.Height() < height {
			height = dj.Height()
		}

		for row := 0; row < height; row++ {
			cmp := strings.Compare(di.Content(row), dj.Content(row))
			switch cmp {
			case -1:
				return true
			case 1:
				return false
			}
		}
		if di.Height() <= dj.Height() {
			return true
		}
		return false
	})

	for _, r := range rows {
		row := tab.Row()
		row.ColumnData(r.key)
		row.ColumnData(r.val)
	}

	return nil
}

func reflectStruct(tab *Tabulate, flags Flags, tags map[string]bool,
	value reflect.Value) error {

loop:
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)

		myFlags := flags
		for _, tag := range strings.Split(field.Tag.Get("tabulate"), ",") {
			if tag == "omitempty" {
				myFlags |= OmitEmpty
			} else if strings.HasPrefix(tag, "@") {
				// Tagged field. Skip unless filter tags contain it.
				if !tags[tag[1:]] {
					continue loop
				}
			}
		}

		v := value.Field(i)

		// Follow pointers.
		for v.Type().Kind() == reflect.Ptr {
			if v.IsZero() {
				if myFlags&OmitEmpty == 0 {
					row := tab.Row()
					row.Column(field.Name)
				}
				continue loop
			}
			v = reflect.Indirect(v)
		}

		if v.CanInterface() {
			switch iv := v.Interface().(type) {
			case encoding.TextMarshaler:
				data, err := iv.MarshalText()
				if err != nil {
					return err
				}
				row := tab.Row()
				row.Column(field.Name)
				row.Column(string(data))
				continue loop
			}
		}

		data, err := reflectValue(tab, flags, tags, v)
		if err != nil {
			return err
		}
		if data.Height() > 0 || flags&OmitEmpty == 0 {
			row := tab.Row()
			row.Column(field.Name)
			row.ColumnData(data)
		}

	}
	return nil
}
