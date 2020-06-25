//
// Copyright (c) 2020 Markku Rossi
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

func ExampleNewUnicode() {
	tab := NewUnicode()

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
	// ┣━━━━━━╋━━━━━━━━╋━━━━━━━━━━━━━━┫
	// ┃ 2018 ┃ 100    ┃ Salary       ┃
	// ┃ 2019 ┃ 110    ┃ Consultation ┃
	// ┃ 2020 ┃ 200    ┃ Lottery      ┃
	// ┗━━━━━━┻━━━━━━━━┻━━━━━━━━━━━━━━┛
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

	tab := NewASCII()
	tab.Header("Key").SetAlign(ML)
	tab.Header("Value")
	err := Reflect(tab, 0, nil, &Book{
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

func ExampleTabulate_MarshalJSON() {
	tab := NewUnicode()
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
