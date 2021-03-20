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

var borderTestBasic = `Year,Income,Expenses
2018,100,90;91;92
2019,110,85
2020,107,50`

var borderTestHdrOnly = `Year,Income,Expenses`

var borderTests = []struct {
	style  Style
	align  Align
	input  string
	result string
}{
	{
		style: Plain,
		align: TL,
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
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
		input: borderTestBasic,
		result: `
        {"2018":["100","90\n91\n92"],"2019":["110","85"],"2020":["107","50"]}
`,
	},
	{
		style: JSON,
		align: MC,
		input: borderTestBasic,
		result: `
        {"2018":["100","90\n91\n92"],"2019":["110","85"],"2020":["107","50"]}
`,
	},
	{
		style: JSON,
		align: BR,
		input: borderTestBasic,
		result: `
        {"2018":["100","90\n91\n92"],"2019":["110","85"],"2020":["107","50"]}
`,
	},

	// Header only tests.
	{
		style: Plain,
		align: TL,
		input: borderTestHdrOnly,
		result: `
Year  Income  Expenses
`,
	},
	{
		style: ASCII,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        +------+--------+----------+
        | Year | Income | Expenses |
        +------+--------+----------+
`,
	},
	{
		style: Unicode,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
        ┃ Year ┃ Income ┃ Expenses ┃
        ┗━━━━━━┻━━━━━━━━┻━━━━━━━━━━┛
`,
	},
	{
		style: UnicodeLight,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        ┌──────┬────────┬──────────┐
        │ Year │ Income │ Expenses │
        └──────┴────────┴──────────┘
`,
	},
	{
		style: UnicodeBold,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓
        ┃ Year ┃ Income ┃ Expenses ┃
        ┗━━━━━━┻━━━━━━━━┻━━━━━━━━━━┛
`,
	},
	{
		style: Colon,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        Year : Income : Expenses
`,
	},
	{
		style: Simple,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        Year Income Expenses
`,
	},
	{
		style: SimpleUnicode,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        Year Income Expenses
`,
	},
	{
		style: SimpleUnicodeBold,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        Year Income Expenses
`,
	},
	{
		style: Github,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        | Year | Income | Expenses |
`,
	},
	{
		style: JSON,
		align: TL,
		input: borderTestHdrOnly,
		result: `
        {}
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
	for idx, test := range borderTests {
		result := tab(test.style, test.align, test.input)
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
