//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func tabulateStyle(style Style) *Tabulate {
	tab := New(style)
	tab.Header("Year").SetAlign(MR)
	tab.Header("Income").SetAlign(MR)

	row := tab.Row()
	row.Column("2018")
	row.Column("100")

	row = tab.Row()
	row.Column("2019")
	row.Column("110")

	row = tab.Row()
	row.Column("2020")
	row.Column("200")

	return tab
}

func ExampleStyle_csv() {
	tabulateStyle(CSV).Print(os.Stdout)
}

func ExampleTabulate_Row() {
	tab := New(Unicode)
	tab.Header("Year").SetAlign(MR)
	tab.Header("Income").SetAlign(MR)

	row := tab.Row()
	row.Column("2018")
	row.Column("100")

	row = tab.Row()
	row.Column("2019")
	row.Column("110")

	row = tab.Row()
	row.Column("2020")
	row.Column("200")

	tab.Print(os.Stdout)
	// Output: ┏━━━━━━┳━━━━━━━━┓
	// ┃ Year ┃ Income ┃
	// ┡━━━━━━╇━━━━━━━━┩
	// │ 2018 │    100 │
	// │ 2019 │    110 │
	// │ 2020 │    200 │
	// └──────┴────────┘
}

func ExampleTabulate_Header() {
	tab := New(Unicode)
	tab.Header("Year").SetAlign(MR)
	tab.Header("Income").SetAlign(MR)
	tab.Print(os.Stdout)
	// Output: ┏━━━━━━┳━━━━━━━━┓
	// ┃ Year ┃ Income ┃
	// ┗━━━━━━┻━━━━━━━━┛
}

func ExampleNew() {
	tab := New(Unicode)

	lines := strings.Split(`Year,Income,Source
2018,100,Salary
2019,110,Consultation
2020,200,Lottery`, "\n")

	// Table headers.
	for _, hdr := range strings.Split(lines[0], ",") {
		tab.Header(hdr)
	}

	// Table data rows.
	for _, line := range lines[1:] {
		row := tab.Row()
		for _, col := range strings.Split(line, ",") {
			row.Column(col)
		}
	}

	tab.Print(os.Stdout)

	// Output: ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━┓
	// ┃ Year ┃ Income ┃ Source       ┃
	// ┡━━━━━━╇━━━━━━━━╇━━━━━━━━━━━━━━┩
	// │ 2018 │ 100    │ Salary       │
	// │ 2019 │ 110    │ Consultation │
	// │ 2020 │ 200    │ Lottery      │
	// └──────┴────────┴──────────────┘
}

func ExampleReflect() {
	type Person struct {
		Name string
	}

	type Book struct {
		Title     string
		Author    []Person
		Publisher string
		Published int
	}

	tab := New(ASCII)
	tab.Header("Key").SetAlign(ML)
	tab.Header("Value")
	err := Reflect(tab, InheritHeaders, nil, &Book{
		Title: "Structure and Interpretation of Computer Programs",
		Author: []Person{
			{
				Name: "Harold Abelson",
			},
			{
				Name: "Gerald Jay Sussman",
			},
			{
				Name: "Julie Sussman",
			},
		},
		Publisher: "MIT Press",
		Published: 1985,
	})
	if err != nil {
		log.Fatal(err)
	}
	tab.Print(os.Stdout)
	// Output: +-----------+---------------------------------------------------+
	// | Key       | Value                                             |
	// +-----------+---------------------------------------------------+
	// | Title     | Structure and Interpretation of Computer Programs |
	// |           | +------+----------------+                         |
	// |           | | Key  | Value          |                         |
	// |           | +------+----------------+                         |
	// |           | | Name | Harold Abelson |                         |
	// |           | +------+----------------+                         |
	// |           | +------+--------------------+                     |
	// |           | | Key  | Value              |                     |
	// | Author    | +------+--------------------+                     |
	// |           | | Name | Gerald Jay Sussman |                     |
	// |           | +------+--------------------+                     |
	// |           | +------+---------------+                          |
	// |           | | Key  | Value         |                          |
	// |           | +------+---------------+                          |
	// |           | | Name | Julie Sussman |                          |
	// |           | +------+---------------+                          |
	// | Publisher | MIT Press                                         |
	// | Published | 1985                                              |
	// +-----------+---------------------------------------------------+
}

func ExampleArray() {
	tab, err := Array(New(ASCII), [][]interface{}{
		{"a", "b", "c"},
		{"1", "2", "3"},
	})
	if err != nil {
		log.Fatal(err)
	}
	tab.Print(os.Stdout)
	// Output: +---+---+---+
	// | a | b | c |
	// +---+---+---+
	// | 1 | 2 | 3 |
	// +---+---+---+
}

func ExampleArray_second() {
	tab, err := Array(New(Unicode), [][]interface{}{
		{"int", "float", "struct"},
		{42, 3.14, struct {
			ival   int
			strval string
		}{
			ival:   42,
			strval: "Hello, world!",
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
	tab.Print(os.Stdout)
	// Output: ┏━━━━━┳━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
	// ┃ int ┃ float ┃ struct                     ┃
	// ┡━━━━━╇━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━━━━┩
	// │ 42  │ 3.14  │ ┌────────┬───────────────┐ │
	// │     │       │ │ ival   │ 42            │ │
	// │     │       │ │ strval │ Hello, world! │ │
	// │     │       │ └────────┴───────────────┘ │
	// └─────┴───────┴────────────────────────────┘
}

func ExampleTabulate_MarshalJSON() {
	tab := New(Unicode)
	tab.Header("Key").SetAlign(MR)
	tab.Header("Value").SetAlign(ML)

	row := tab.Row()
	row.Column("Boolean")
	row.ColumnData(NewValue(false))

	row = tab.Row()
	row.Column("Integer")
	row.ColumnData(NewValue(42))

	data, err := json.Marshal(tab)
	if err != nil {
		log.Fatalf("JSON marshal failed: %s", err)
	}
	fmt.Println(string(data))
	// Output: {"Boolean":false,"Integer":42}
}
