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
	return reflectValue(tab, "", reflect.ValueOf(v))
}

func reflectValue(tab *Tabulate, name string, value reflect.Value) error {
	var text string

	switch value.Type().Kind() {
	case reflect.Bool:
		text = fmt.Sprintf("%v", value.Bool())

	case reflect.Int:
		text = fmt.Sprintf("%v", value.Int())

	case reflect.Ptr:
		return reflectValue(tab, name, reflect.Indirect(value))

	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			// XXX Tags

			v := value.Field(i)
			var err error

			switch v.Type().Kind() {
			case reflect.Ptr:
				if v.IsZero() {
				} else if reflect.Indirect(v).Type().Kind() == reflect.Struct {
					sub := tab.Clone()
					err = reflectValue(sub, value.Type().Field(i).Name, v)
					if err != nil {
						return err
					}
					row := tab.Row()
					row.Column(NewText(value.Type().Field(i).Name))
					row.Column(sub.Data())
				} else {
					err = reflectValue(tab, value.Type().Field(i).Name, v)
				}

			case reflect.Struct:
				sub := tab.Clone()
				err = reflectValue(sub, value.Type().Field(i).Name, v)
				if err != nil {
					return err
				}
				row := tab.Row()
				row.Column(NewText(value.Type().Field(i).Name))
				row.Column(sub.Data())

			default:
				err := reflectValue(tab, value.Type().Field(i).Name, v)
				if err != nil {
					return err
				}
			}
		}
		return nil

	default:
		text = value.String()
	}

	row := tab.Row()
	row.Column(NewText(name))
	row.Column(NewText(text))

	return nil
}
