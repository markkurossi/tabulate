//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"fmt"
	"reflect"
)

func Reflect(tab *Tabulate, v interface{}) error {
	value := reflect.ValueOf(v)

	// Follows pointers.
	for value.Type().Kind() == reflect.Ptr {
		if value.IsZero() {
			return nil
		}
		value = reflect.Indirect(value)
	}

	if value.Type().Kind() == reflect.Struct {
		return reflectStruct(tab, value)
	}

	lines, err := reflectValue(tab, value)
	if err != nil {
		return err
	}
	row := tab.Row()
	row.Column("")
	row.ColumnData(NewLinesData(lines))

	return nil
}

func reflectValue(tab *Tabulate, value reflect.Value) ([]string, error) {
	var text string

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
					lines = append(lines, "<nil>")
					continue loop
				}
			}
			switch v.Type().Kind() {
			case reflect.Struct:
				sub := tab.Clone()
				err := reflectStruct(sub, v)
				if err != nil {
					return nil, err
				}
				data := sub.Data()
				for row := 0; row < data.Height(); row++ {
					lines = append(lines, data.Content(row))
				}

			default:
				l, err := reflectValue(tab, v)
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

	return []string{text}, nil
}

func reflectStruct(tab *Tabulate, value reflect.Value) error {
	for i := 0; i < value.NumField(); i++ {
		// XXX Tags

		v := value.Field(i)

		// Follow pointers.
		for v.Type().Kind() == reflect.Ptr {
			if v.IsZero() {
				continue
			}
			v = reflect.Indirect(v)
		}

		var err error

		switch v.Type().Kind() {
		case reflect.Struct:
			sub := tab.Clone()
			err = reflectStruct(sub, v)
			if err != nil {
				return err
			}
			row := tab.Row()
			row.Column(value.Type().Field(i).Name)
			row.ColumnData(sub.Data())

		default:
			lines, err := reflectValue(tab, v)
			if err != nil {
				return err
			}
			row := tab.Row()
			row.Column(value.Type().Field(i).Name)
			row.ColumnData(NewLinesData(lines))
		}
	}
	return nil
}
