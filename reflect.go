//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

// Flags control how reflection tabulation operates on different
// values.
type Flags int

// Flag values for reflection tabulation.
const (
	OmitEmpty Flags = 1 << iota
)

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

	lines, err := reflectValue(tab, flags, tagMap, value)
	if err != nil {
		return err
	}
	row := tab.Row()
	row.Column("")
	row.ColumnData(NewLinesData(lines))

	return nil
}

func reflectValue(tab *Tabulate, flags Flags, tags map[string]bool,
	value reflect.Value) ([]string, error) {
	var text string

	if value.CanInterface() {
		switch v := value.Interface().(type) {
		case encoding.TextMarshaler:
			data, err := v.MarshalText()
			if err != nil {
				return nil, err
			}
			return []string{string(data)}, nil
		}
	}

	switch value.Type().Kind() {
	case reflect.Bool:
		text = fmt.Sprintf("%v", value.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		text = fmt.Sprintf("%v", value.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		text = fmt.Sprintf("%v", value.Uint())

	case reflect.Slice:
		var lines []string
	loop:
		for i := 0; i < value.Len(); i++ {
			v := value.Index(i)
			// Follow pointers.
			for v.Type().Kind() == reflect.Ptr {
				if v.IsZero() {
					if flags&OmitEmpty == 0 {
						lines = append(lines, "<nil>")
					}
					continue loop
				}
				v = reflect.Indirect(v)
			}
			switch v.Type().Kind() {
			case reflect.Struct:
				sub := tab.Clone()
				err := reflectStruct(sub, flags, tags, v)
				if err != nil {
					return nil, err
				}
				data := sub.Data()
				for row := 0; row < data.Height(); row++ {
					lines = append(lines, data.Content(row))
				}

			default:
				l, err := reflectValue(tab, flags, tags, v)
				if err != nil {
					return nil, err
				}
				lines = append(lines, l...)
			}
		}
		return lines, nil

	default:
		text = value.String()
	}

	if len(text) == 0 && flags&OmitEmpty == 1 {
		return nil, nil
	}

	return []string{text}, nil
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

		var err error

		switch v.Type().Kind() {
		case reflect.Struct:
			sub := tab.Clone()
			err = reflectStruct(sub, flags, tags, v)
			if err != nil {
				return err
			}
			row := tab.Row()
			row.Column(field.Name)
			row.ColumnData(sub.Data())

		case reflect.Map:
			iter := v.MapRange()
			sub := tab.Clone()
			for iter.Next() {
				row := sub.Row()
				lines, err := reflectValue(tab, flags, tags, iter.Key())
				if err != nil {
					return err
				}
				row.ColumnData(NewLinesData(lines))
				lines, err = reflectValue(tab, flags, tags, iter.Value())
				if err != nil {
					return err
				}
				row.ColumnData(NewLinesData(lines))
			}
			row := tab.Row()
			row.Column(field.Name)
			row.ColumnData(sub.Data())

		default:
			lines, err := reflectValue(tab, flags, tags, v)
			if err != nil {
				return err
			}
			if len(lines) > 0 || flags&OmitEmpty == 0 {
				row := tab.Row()
				row.Column(field.Name)
				row.ColumnData(NewLinesData(lines))
			}
		}

	}
	return nil
}
