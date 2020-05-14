//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
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
			Person{
				Name: "Harold Abelson",
			},
			Person{
				Name: "Gerald Jay Sussman",
			},
			Person{
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
