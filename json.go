//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"encoding/json"
	"errors"
	"fmt"
)

type jsonMarshaler interface {
	marshalJSON() (map[string]interface{}, error)
}

// MarshalJSON implements the JSON Marshaler interface.
func (t *Tabulate) MarshalJSON() ([]byte, error) {
	content, err := t.marshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(content)
}

func (t *Tabulate) marshalJSON() (map[string]interface{}, error) {
	content := make(map[string]interface{})

	for _, row := range t.Rows {
		if len(row.Columns) < 2 {
			return nil, errors.New("JSON tabulation must have at least two columns")
		}
		var columns []interface{}
		for i := 1; i < len(row.Columns); i++ {
			col := row.Columns[i]
			fmt.Printf("Data: %v (%T)\n", col.Data, col.Data)
			marshaler, ok := col.Data.(jsonMarshaler)
			if ok {
				v, err := marshaler.marshalJSON()
				if err != nil {
					return nil, err
				}
				columns = append(columns, v)
			} else {
				columns = append(columns, col.Data.String())
			}
		}
		key := row.Columns[0].Data.String()
		if len(columns) > 1 {
			content[key] = columns
		} else {
			content[key] = columns[0]
		}
	}

	return content, nil
}
