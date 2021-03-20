//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"os"
	"strings"
	"testing"
)

var borderTests = []struct {
	style  Style
	align  Align
	result string
}{
	{
		style: Plain,
		align: TL,
		result: `
Year  Income  Expenses
2018  100     90
              91
              92
2019  110     85
2020  107     50
`,
	},
	{
		style: Plain,
		align: MC,
		result: `
Year  Income  Expenses
                 90
2018   100       91
                 92
2019   110       85
2020   107       50
`,
	},
	{
		style: Plain,
		align: BR,
		result: `
Year  Income  Expenses
                    90
                    91
2018     100        92
2019     110        85
2020     107        50
`,
	},
	{
		style: ASCII,
		align: TL,
		result: `
+------+--------+----------+
| Year | Income | Expenses |
+------+--------+----------+
| 2018 | 100    | 90       |
|      |        | 91       |
|      |        | 92       |
| 2019 | 110    | 85       |
| 2020 | 107    | 50       |
+------+--------+----------+
`,
	},
	{
		style: ASCII,
		align: MC,
		result: `
+------+--------+----------+
| Year | Income | Expenses |
+------+--------+----------+
|      |        |    90    |
| 2018 |  100   |    91    |
|      |        |    92    |
| 2019 |  110   |    85    |
| 2020 |  107   |    50    |
+------+--------+----------+
`,
	},
	{
		style: ASCII,
		align: BR,
		result: `
+------+--------+----------+
| Year | Income | Expenses |
+------+--------+----------+
|      |        |       90 |
|      |        |       91 |
| 2018 |    100 |       92 |
| 2019 |    110 |       85 |
| 2020 |    107 |       50 |
+------+--------+----------+
`,
	},
	{
		style: Unicode,
		align: TL,
		result: `
┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
┃ Year ┃ Income ┃ Expenses ┃
┡━━━━━━╇━━━━━━━━╇━━━━━━━━━━┩
│ 2018 │ 100    │ 90       │
│      │        │ 91       │
│      │        │ 92       │
│ 2019 │ 110    │ 85       │
│ 2020 │ 107    │ 50       │
└──────┴────────┴──────────┘
`,
	},
	{
		style: Unicode,
		align: MC,
		result: `
┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
┃ Year ┃ Income ┃ Expenses ┃
┡━━━━━━╇━━━━━━━━╇━━━━━━━━━━┩
│      │        │    90    │
│ 2018 │  100   │    91    │
│      │        │    92    │
│ 2019 │  110   │    85    │
│ 2020 │  107   │    50    │
└──────┴────────┴──────────┘
`,
	},
	{
		style: Unicode,
		align: BR,
		result: `
┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
┃ Year ┃ Income ┃ Expenses ┃
┡━━━━━━╇━━━━━━━━╇━━━━━━━━━━┩
│      │        │       90 │
│      │        │       91 │
│ 2018 │    100 │       92 │
│ 2019 │    110 │       85 │
│ 2020 │    107 │       50 │
└──────┴────────┴──────────┘
`,
	},
	{
		style: UnicodeLight,
		align: TL,
		result: `
        ┌──────┬────────┬──────────┐
        │ Year │ Income │ Expenses │
        ├──────┼────────┼──────────┤
        │ 2018 │ 100    │ 90       │
        │      │        │ 91       │
        │      │        │ 92       │
        │ 2019 │ 110    │ 85       │
        │ 2020 │ 107    │ 50       │
        └──────┴────────┴──────────┘
`,
	},
	{
		style: UnicodeLight,
		align: MC,
		result: `
        ┌──────┬────────┬──────────┐
        │ Year │ Income │ Expenses │
        ├──────┼────────┼──────────┤
        │      │        │    90    │
        │ 2018 │  100   │    91    │
        │      │        │    92    │
        │ 2019 │  110   │    85    │
        │ 2020 │  107   │    50    │
        └──────┴────────┴──────────┘
`,
	},
	{
		style: UnicodeLight,
		align: BR,
		result: `
        ┌──────┬────────┬──────────┐
        │ Year │ Income │ Expenses │
        ├──────┼────────┼──────────┤
        │      │        │       90 │
        │      │        │       91 │
        │ 2018 │    100 │       92 │
        │ 2019 │    110 │       85 │
        │ 2020 │    107 │       50 │
        └──────┴────────┴──────────┘
`,
	},
	{
		style: UnicodeBold,
		align: TL,
		result: `
        ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
        ┃ Year ┃ Income ┃ Expenses ┃
        ┣━━━━━━╋━━━━━━━━╋━━━━━━━━━━┫
        ┃ 2018 ┃ 100    ┃ 90       ┃
        ┃      ┃        ┃ 91       ┃
        ┃      ┃        ┃ 92       ┃
        ┃ 2019 ┃ 110    ┃ 85       ┃
        ┃ 2020 ┃ 107    ┃ 50       ┃
        ┗━━━━━━┻━━━━━━━━┻━━━━━━━━━━┛
`,
	},
	{
		style: UnicodeBold,
		align: MC,
		result: `
        ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
        ┃ Year ┃ Income ┃ Expenses ┃
        ┣━━━━━━╋━━━━━━━━╋━━━━━━━━━━┫
        ┃      ┃        ┃    90    ┃
        ┃ 2018 ┃  100   ┃    91    ┃
        ┃      ┃        ┃    92    ┃
        ┃ 2019 ┃  110   ┃    85    ┃
        ┃ 2020 ┃  107   ┃    50    ┃
        ┗━━━━━━┻━━━━━━━━┻━━━━━━━━━━┛
`,
	},
	{
		style: UnicodeBold,
		align: BR,
		result: `
        ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
        ┃ Year ┃ Income ┃ Expenses ┃
        ┣━━━━━━╋━━━━━━━━╋━━━━━━━━━━┫
        ┃      ┃        ┃       90 ┃
        ┃      ┃        ┃       91 ┃
        ┃ 2018 ┃    100 ┃       92 ┃
        ┃ 2019 ┃    110 ┃       85 ┃
        ┃ 2020 ┃    107 ┃       50 ┃
        ┗━━━━━━┻━━━━━━━━┻━━━━━━━━━━┛
`,
	},
	{
		style: Colon,
		align: TL,
		result: `
        Year : Income : Expenses
        2018 : 100    : 90
             :        : 91
             :        : 92
        2019 : 110    : 85
        2020 : 107    : 50
`,
	},
	{
		style: Colon,
		align: MC,
		result: `
        Year : Income : Expenses
             :        :    90
        2018 :  100   :    91
             :        :    92
        2019 :  110   :    85
        2020 :  107   :    50
`,
	},
	{
		style: Colon,
		align: BR,
		result: `
        Year : Income : Expenses
             :        :       90
             :        :       91
        2018 :    100 :       92
        2019 :    110 :       85
        2020 :    107 :       50
`,
	},
	{
		style: Simple,
		align: TL,
		result: `
        Year Income Expenses
        ---- ------ --------
        2018 100    90
                    91
                    92
        2019 110    85
        2020 107    50
`,
	},
	{
		style: Simple,
		align: MC,
		result: `
        Year Income Expenses
        ---- ------ --------
                       90
        2018  100      91
                       92
        2019  110      85
        2020  107      50
`,
	},
	{
		style: Simple,
		align: BR,
		result: `
        Year Income Expenses
        ---- ------ --------
                          90
                          91
        2018    100       92
        2019    110       85
        2020    107       50
`,
	},
	{
		style: SimpleUnicode,
		align: TL,
		result: `
        Year Income Expenses
        ──── ────── ────────
        2018 100    90
                    91
                    92
        2019 110    85
        2020 107    50
`,
	},
	{
		style: SimpleUnicode,
		align: MC,
		result: `
        Year Income Expenses
        ──── ────── ────────
                       90
        2018  100      91
                       92
        2019  110      85
        2020  107      50
`,
	},
	{
		style: SimpleUnicode,
		align: BR,
		result: `
        Year Income Expenses
        ──── ────── ────────
                          90
                          91
        2018    100       92
        2019    110       85
        2020    107       50
`,
	},
	{
		style: SimpleUnicodeBold,
		align: TL,
		result: `
        Year Income Expenses
        ━━━━ ━━━━━━ ━━━━━━━━
        2018 100    90
                    91
                    92
        2019 110    85
        2020 107    50
`,
	},
	{
		style: SimpleUnicodeBold,
		align: MC,
		result: `
        Year Income Expenses
        ━━━━ ━━━━━━ ━━━━━━━━
                       90
        2018  100      91
                       92
        2019  110      85
        2020  107      50
`,
	},
	{
		style: SimpleUnicodeBold,
		align: BR,
		result: `
        Year Income Expenses
        ━━━━ ━━━━━━ ━━━━━━━━
                          90
                          91
        2018    100       92
        2019    110       85
        2020    107       50
`,
	},
	{
		style: Github,
		align: TL,
		result: `
        | Year | Income | Expenses |
        |------|--------|----------|
        | 2018 | 100    | 90       |
        |      |        | 91       |
        |      |        | 92       |
        | 2019 | 110    | 85       |
        | 2020 | 107    | 50       |
`,
	},
	{
		style: Github,
		align: MC,
		result: `
        | Year | Income | Expenses |
        |------|--------|----------|
        |      |        |    90    |
        | 2018 |  100   |    91    |
        |      |        |    92    |
        | 2019 |  110   |    85    |
        | 2020 |  107   |    50    |
`,
	},
	{
		style: Github,
		align: BR,
		result: `
        | Year | Income | Expenses |
        |------|--------|----------|
        |      |        |       90 |
        |      |        |       91 |
        | 2018 |    100 |       92 |
        | 2019 |    110 |       85 |
        | 2020 |    107 |       50 |
`,
	},
	{
		style: CSV,
		align: TL,
		result: `
        Year,Income,Expenses
        2018,100,90
        ,,91
        ,,92
        2019,110,85
        2020,107,50
`,
	},
	{
		style: CSV,
		align: MC,
		result: `
        Year,Income,Expenses
        ,,90
        2018,100,91
        ,,92
        2019,110,85
        2020,107,50
`,
	},
	{
		style: CSV,
		align: BR,
		result: `
        Year,Income,Expenses
        ,,90
        ,,91
        2018,100,92
        2019,110,85
        2020,107,50
`,
	},
	{
		style: JSON,
		align: TL,
		result: `
        {"2018":["100","90\n91\n92"],"2019":["110","85"],"2020":["107","50"]}
`,
	},
	{
		style: JSON,
		align: MC,
		result: `
        {"2018":["100","90\n91\n92"],"2019":["110","85"],"2020":["107","50"]}
`,
	},
	{
		style: JSON,
		align: BR,
		result: `
        {"2018":["100","90\n91\n92"],"2019":["110","85"],"2020":["107","50"]}
`,
	},
}

func tab(style Style, align Align, data string) string {
	tab := New(style)

	rows := strings.Split(data, "\n")
	for _, hdr := range strings.Split(rows[0], ",") {
		tab.Header(hdr).SetAlign(align)
	}

	for _, r := range rows[1:] {
		row := tab.Row()
		for _, col := range strings.Split(r, ",") {
			row.ColumnData(NewLinesData(strings.Split(col, ";")))
		}
	}
	var sb strings.Builder
	tab.Print(&sb)
	return sb.String()
}

func cleanup(input string) []string {
	var result []string

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			result = append(result, line)
		}
	}
	return result
}

func match(a, b string) bool {
	aLines := cleanup(a)
	bLines := cleanup(b)

	if len(aLines) != len(bLines) {
		return false
	}
	for idx, al := range aLines {
		if al != bLines[idx] {
			return false
		}
	}
	return true
}

func TestStyles(t *testing.T) {
	data := `Year,Income,Expenses
2018,100,90;91;92
2019,110,85
2020,107,50`

	for idx, test := range borderTests {
		result := tab(test.style, test.align, data)
		if !match(result, test.result) {
			t.Errorf("TestStyles %d: got:\n%s\nexpected:\n%s\n", idx,
				result, test.result)
		}
	}
}

func tabulateRows(tab *Tabulate, align Align, rows []string) *Tabulate {

	if len(rows[0]) > 0 {
		for _, hdr := range strings.Split(rows[0], ",") {
			tab.Header(hdr).SetAlign(align)
		}
	}

	for i := 1; i < len(rows); i++ {
		row := tab.Row()
		for _, col := range strings.Split(rows[i], ",") {
			row.ColumnData(NewText(col))
		}
	}
	return tab
}

func tabulate(tab *Tabulate, align Align, data string) *Tabulate {
	if len(data) == 0 {
		return tab
	}
	return tabulateRows(tab, align, strings.Split(data, "\n"))
}

func align(align Align, data string) {
	tabulate(New(Plain), align, data).Print(os.Stdout)
	tabulate(New(ASCII), align, data).Print(os.Stdout)
	tabulate(New(Unicode), align, data).Print(os.Stdout)
	tabulate(New(UnicodeLight), align, data).Print(os.Stdout)
	tabulate(New(UnicodeBold), align, data).Print(os.Stdout)
	tabulate(New(Colon), align, data).Print(os.Stdout)
	tabulate(New(Simple), align, data).Print(os.Stdout)
	tabulate(New(SimpleUnicode), align, data).Print(os.Stdout)
	tabulate(New(SimpleUnicodeBold), align, data).Print(os.Stdout)
	tabulate(New(Github), align, data).Print(os.Stdout)
	tabulate(New(JSON), align, data).Print(os.Stdout)
}

func TestEmpty(t *testing.T) {
	data := `Year,Income,Expenses`
	align(TL, data)
}

func TestRowsOnly(t *testing.T) {
	data := `
2018,100,90
2019,110,85
2020,107,50`

	align(TL, data)
}

func TestNoColumns(t *testing.T) {
	align(TL, "")
}

var csv = `Year,Income,Source|2018,100,Salary|2019,110,"Consultation"|2020,120,Lottery
et al`

func TestCSV(t *testing.T) {
	rows := strings.Split(csv, "|")
	tabulateRows(New(CSV), None, rows).Print(os.Stdout)
}

var missingCols = `Year,Value
2018,100
2019,
2020,100,200`

func TestMissingColumns(t *testing.T) {
	rows := strings.Split(missingCols, "\n")
	tabulateRows(New(Unicode), TL, rows).Print(os.Stdout)
}

func TestNested(t *testing.T) {
	tab := New(Unicode)

	tab.Header("Key").SetAlign(MR)
	tab.Header("Value").SetAlign(MC)

	row := tab.Row()
	row.Column("Name")
	row.Column("ACME Corp.")

	row = tab.Row()
	row.Column("Numbers")

	data := `Year,Income,Expenses
2018,100,90
2019,110,85
2020,107,50`

	row.ColumnData(tabulate(New(Unicode), TR, data))

	tab.Print(os.Stdout)
}

func TestWide(t *testing.T) {
	tab := New(ASCII)

	tab.Header("我")
	tab.Header("是")
	tab.Header("测试")

	row := tab.Row()
	row.Column("a")
	row.Column("a")
	row.Column("a")

	tab.Print(os.Stdout)
}
