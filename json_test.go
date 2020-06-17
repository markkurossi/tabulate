//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSONTimeSeries(t *testing.T) {
	tab := tabulate(NewWS(), TL)
	data, err := json.MarshalIndent(tab, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal time series failed: %s", err)
	}
	fmt.Printf("JSON time series:\n%s\n", string(data))
}

func TestJSONReflect(t *testing.T) {
	tab := NewWS()
	tab.Header("Field")
	tab.Header("Value")

	err := Reflect(tab, OmitEmpty, nil, &Outer{
		Name: "Alyssa P. Hacker",
		Age:  45,
		Address: &Address{
			Lines: []string{"42 Hacker way", "03139 Cambridge", "MA"},
		},
		Info: []*Info{
			{
				Email: "mtr@iki.fi",
			},
			{
				Email: "markku.rossi@gmail.com",
				Work:  true,
			},
		},
		Mapping: map[string]string{
			"First":  "1st",
			"Second": "2nd",
		},
	})
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
	data, err := json.MarshalIndent(tab, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal reflect failed: %s", err)
	}
	fmt.Printf("JSON marshal reflect:\n%s\n", string(data))
}
