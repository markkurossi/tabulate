//
// Copyright (c) 2020-2025 Markku Rossi
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
	rows := `Year,Income,Expenses
2018,100,90
2019,110,85
2020,107,50`

	tab := tabulate(New(Plain), TL, rows)
	data, err := json.MarshalIndent(tab, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal time series failed: %s", err)
	}
	expected := `
        {
          "2018": [
            "100",
            "90"
          ],
          "2019": [
            "110",
            "85"
          ],
          "2020": [
            "107",
            "50"
          ]
        }
`

	match(t, string(data), expected, "TestJSONTimeSeries")
}

func TestJSONReflect(t *testing.T) {
	tab := New(Plain)
	tab.Header("Field")
	tab.Header("Value")

	err := Reflect(tab, OmitEmpty, nil, &Outer{
		Name: "Alyssa P. Hacker",
		Age:  45,
		NPS:  9.9,
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
		Meta: &Info{
			Email: "mtr@iki.fi",
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
	expected := `
        {
          "Address": {
            "Lines": [
              "42 Hacker way",
              "03139 Cambridge",
              "MA"
            ]
          },
          "Age": 45,
          "Info": [
            {
              "Email": "mtr@iki.fi",
              "Work": false
            },
            {
              "Email": "markku.rossi@gmail.com",
              "Work": true
            }
          ],
          "Mapping": {
            "First": "1st",
            "Second": "2nd"
          },
          "Meta": {
            "Email": "mtr@iki.fi",
            "Work": false
          },
          "NPS": 9.9,
          "Name": "Alyssa P. Hacker"
        }
`

	match(t, string(data), expected, "TestJSONReflect")
}

func XTestJSONCertReflect(t *testing.T) {
	c, err := decodeCertificate()
	if err != nil {
		t.Fatalf("Failed to decode certificate: %s", err)
	}
	tab := New(Plain)
	tab.Header("Field")
	tab.Header("Value")

	err = Reflect(tab, OmitEmpty, nil, c)
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
	data, err := json.MarshalIndent(tab, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal cert reflect failed: %s", err)
	}
	if false {
		fmt.Printf("JSON marshal cert reflect:\n%s\n", string(data))
	}
}
