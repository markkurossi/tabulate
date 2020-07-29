# tabulate

Tabulate is an utility library for making simple data
visualizations. Tabulate works on tabular data. The data tables can be
constructed explicity by calling the row and column functions, or with
reflection from Go values.

[![Build Status](https://img.shields.io/github/workflow/status/markkurossi/tabulate/Go)](https://github.com/markkurossi/tabulate/actions)
[![Git Hub](https://img.shields.io/github/last-commit/markkurossi/tabulate.svg)](https://github.com/markkurossi/tabulate/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/markkurossi/tabulate)](https://goreportcard.com/report/github.com/markkurossi/tabulate)

## Programmatic table construction

In the programmatic table construction, you first create a new table
and define the headers with optional layout attributes:

```go
tab := tabulate.New(tabulate.Unicode)
tab.Header("Year").SetAlign(tabulate.MR)
tab.Header("Income").SetAlign(tabulate.MR)
```

After that, you add data rows:

```go
row := tab.Row()
row.Column("2018")
row.Column("100")

row = tab.Row()
row.Column("2019")
row.Column("110")

row = tab.Row()
row.Column("2020")
row.Column("200")
```

Finally, you print the table:

```go
tab.Print(os.Stdout)
```

This outputs the table to the selected writer:

    ┏━━━━━━┳━━━━━━━━┓
    ┃ Year ┃ Income ┃
    ┡━━━━━━╇━━━━━━━━┩
    │ 2018 │    100 │
    │ 2019 │    110 │
    │ 2020 │    200 │
    └──────┴────────┘

## Reflection

The reflection mode allows you to easily tabulate Go data
structures. The resulting table will always have two columns: key and
value. But the value columns can contain nested tables.

```go
type Person struct {
    Name string
}

type Book struct {
    Title     string
    Author    []Person
    Publisher string
    Published int
}

tab := tabulate.New(tabulate.ASCII)
tab.Header("Key").SetAlign(tabulate.ML)
tab.Header("Value")
err := tabulate.Reflect(tab, 0, nil, &Book{
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
```

This example renders the following table:

    +-----------+---------------------------------------------------+
    | Key       | Value                                             |
    +-----------+---------------------------------------------------+
    | Title     | Structure and Interpretation of Computer Programs |
    |           | +------+----------------+                         |
    |           | | Key  | Value          |                         |
    |           | +------+----------------+                         |
    |           | | Name | Harold Abelson |                         |
    |           | +------+----------------+                         |
    |           | +------+--------------------+                     |
    |           | | Key  | Value              |                     |
    | Author    | +------+--------------------+                     |
    |           | | Name | Gerald Jay Sussman |                     |
    |           | +------+--------------------+                     |
    |           | +------+---------------+                          |
    |           | | Key  | Value         |                          |
    |           | +------+---------------+                          |
    |           | | Name | Julie Sussman |                          |
    |           | +------+---------------+                          |
    | Publisher | MIT Press                                         |
    | Published | 1985                                              |
    +-----------+---------------------------------------------------+

# Formatting

## Cell alignment

Column headers set the default alignment for all cells in the
corresponding columns. The column default alignment is set when the
headers are defined:

```go
tab.Header("Year").SetAlign(tabulate.MR)
```

The alignment is defined with the Align constants. The first character
of the constant name specifies the vertical alignment (Top, Middle,
Bottom) and the second character specifies the horizointal alignment
(Left, Center, Right).

| Alignment | Vertical | Horizontal |
|:---------:|:--------:|:----------:|
| TL        | Top      | Left       |
| TC        | Top      | Center     |
| TR        | Top      | Right      |
| ML        | Middle   | Left       |
| MC        | Middle   | Center     |
| MR        | Middle   | Right      |
| BL        | Bottom   | Left       |
| BC        | Bottom   | Center     |
| BR        | Bottom   | Right      |
| None      | -        | -          |

The default alignment can be overridden by calling the SetAlign() for
the data column:

```go
row = tab.Row()
row.Column("Integer").SetAlign(tabulate.TL)
```

# Output formats

## Plain

The Plain format does not draw any table borders:

     Year  Income  Expenses
     2018  100     90
     2019  110     85
     2020  107     50

## ASCII

The ASCII format creates a new tabulator that uses ASCII characters to
render the table borders:

    +------+--------+----------+
    | Year | Income | Expenses |
    +------+--------+----------+
    | 2018 | 100    | 90       |
    | 2019 | 110    | 85       |
    | 2020 | 107    | 50       |
    +------+--------+----------+

## Unicode

The Unicode format creates a new tabulator that uses Unicode line
drawing characters to render the table borders:

    ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
    ┃ Year ┃ Income ┃ Expenses ┃
    ┡━━━━━━╇━━━━━━━━╇━━━━━━━━━━┩
    │ 2018 │ 100    │ 90       │
    │ 2019 │ 110    │ 85       │
    │ 2020 │ 107    │ 50       │
    └──────┴────────┴──────────┘

## UnicodeLight

The UnicodeLight format creates a new tabulator that uses thin Unicode
line drawing characters to render the table borders:

    ┌──────┬────────┬──────────┐
    │ Year │ Income │ Expenses │
    ├──────┼────────┼──────────┤
    │ 2018 │ 100    │ 90       │
    │ 2019 │ 110    │ 85       │
    │ 2020 │ 107    │ 50       │
    └──────┴────────┴──────────┘

## UnicodeBold

The UnicodeBold format creates a new tabulator that uses thick Unicode
line drawing characters to render the table borders:

    ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
    ┃ Year ┃ Income ┃ Expenses ┃
    ┣━━━━━━╋━━━━━━━━╋━━━━━━━━━━┫
    ┃ 2018 ┃ 100    ┃ 90       ┃
    ┃ 2019 ┃ 110    ┃ 85       ┃
    ┃ 2020 ┃ 107    ┃ 50       ┃
    ┗━━━━━━┻━━━━━━━━┻━━━━━━━━━━┛

## Colon

The Colon format creates a new tabulator that uses colon (':')
character to render vertical table borders:

    Year : Income : Expenses
    2018 : 100    : 90
    2019 : 110    : 85
    2020 : 107    : 50

## Simple

The Simple format draws horizontal lines between header and body
columns:

    Year  Income  Expenses
    ----  ------  --------
    2018  100     90
    2019  110     85
    2020  107     50

## Github

The Github creates tables with the Github Markdown syntax:

    | Year | Income | Expenses |
    |------|--------|----------|
    | 2018 | 100    | 90       |
    | 2019 | 110    | 85       |
    | 2020 | 107    | 50       |

## Comma-Separated Values (CSV) output

The NewCSV() creates a new tabulator that outputs the data in CSV
format. It uses empty borders and an escape function which handles ','
and '\n' characters inside cell values:

    Year,Income,Source
    2018,100,Salary
    2019,110,"""Consultation"""
    2020,120,"Lottery
    et al"

## JSON output

The NewJSON() creates a new tabulator that outputs the data in JSON
format:

    {"2018":["100","90"],"2019":["110","85"],"2020":["107","50"]}

## Native JSON marshalling

The Tabulate object implements the MarshalJSON interface so you can
marshal tabulated data directly into JSON.


```go
tab := tabulate.New(tabulate.Unicode)
tab.Header("Key").SetAlign(tabulate.MR)
tab.Header("Value").SetAlign(tabulate.ML)

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
```

This example outputs the following JSON output:

```json
{"Boolean":false,"Integer":42}
```
