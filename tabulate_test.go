//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"fmt"
	"strings"
	"testing"
)

var borderTestBasic = `Year,Income,Expenses
2018,100,90;91;92
2019,110,85
2020,107,50`

var borderTestHdrOnly = `Year,Income,Expenses`

var borderTestBodyOnly = `
2018,100,9000
2019,110,85;86;86
2020,107,50`

var borderTests = []struct {
	style  Style
	align  Align
	input  string
	rowSep string
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
		style: CompactUnicode,
		align: TL,
		input: borderTestBasic,
		result: `
        ┏━━━━┳━━━━━━┳━━━━━━━━┓
        ┃Year┃Income┃Expenses┃
        ┡━━━━╇━━━━━━╇━━━━━━━━┩
        │2018│100   │90      │
        │    │      │91      │
        │    │      │92      │
        │2019│110   │85      │
        │2020│107   │50      │
        └────┴──────┴────────┘
`,
	},
	{
		style: CompactUnicodeLight,
		align: TL,
		input: borderTestBasic,
		result: `
        ┌────┬──────┬────────┐
        │Year│Income│Expenses│
        ├────┼──────┼────────┤
        │2018│100   │90      │
        │    │      │91      │
        │    │      │92      │
        │2019│110   │85      │
        │2020│107   │50      │
        └────┴──────┴────────┘
`,
	},
	{
		style: CompactUnicodeBold,
		align: TL,
		input: borderTestBasic,
		result: `
        ┏━━━━┳━━━━━━┳━━━━━━━━┓
        ┃Year┃Income┃Expenses┃
        ┣━━━━╋━━━━━━╋━━━━━━━━┫
        ┃2018┃100   ┃90      ┃
        ┃    ┃      ┃91      ┃
        ┃    ┃      ┃92      ┃
        ┃2019┃110   ┃85      ┃
        ┃2020┃107   ┃50      ┃
        ┗━━━━┻━━━━━━┻━━━━━━━━┛
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

	// Body only.
	{
		style: Plain,
		align: TL,
		input: borderTestBodyOnly,
		result: `
         2018  100  9000
         2019  110  85
                    86
                    86
         2020  107  50
`,
	},
	{
		style: Plain,
		align: MC,
		input: borderTestBodyOnly,
		result: `
         2018  100  9000
                     85
         2019  110   86
                     86
         2020  107   50
`,
	},
	{
		style: Plain,
		align: BR,
		input: borderTestBodyOnly,
		result: `
         2018  100  9000
                      85
                      86
         2019  110    86
         2020  107    50
`,
	},
	{
		style: ASCII,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        +------+-----+------+
        | 2018 | 100 | 9000 |
        | 2019 | 110 | 85   |
        |      |     | 86   |
        |      |     | 86   |
        | 2020 | 107 | 50   |
        +------+-----+------+
`,
	},
	{
		style: ASCII,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        +------+-----+------+
        | 2018 | 100 | 9000 |
        |      |     |  85  |
        | 2019 | 110 |  86  |
        |      |     |  86  |
        | 2020 | 107 |  50  |
        +------+-----+------+
`,
	},
	{
		style: ASCII,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        +------+-----+------+
        | 2018 | 100 | 9000 |
        |      |     |   85 |
        |      |     |   86 |
        | 2019 | 110 |   86 |
        | 2020 | 107 |   50 |
        +------+-----+------+
`,
	},
	{
		style: Unicode,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        ┌──────┬─────┬──────┐
        │ 2018 │ 100 │ 9000 │
        │ 2019 │ 110 │ 85   │
        │      │     │ 86   │
        │      │     │ 86   │
        │ 2020 │ 107 │ 50   │
        └──────┴─────┴──────┘
`,
	},
	{
		style: Unicode,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        ┌──────┬─────┬──────┐
        │ 2018 │ 100 │ 9000 │
        │      │     │  85  │
        │ 2019 │ 110 │  86  │
        │      │     │  86  │
        │ 2020 │ 107 │  50  │
        └──────┴─────┴──────┘
`,
	},
	{
		style: Unicode,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        ┌──────┬─────┬──────┐
        │ 2018 │ 100 │ 9000 │
        │      │     │   85 │
        │      │     │   86 │
        │ 2019 │ 110 │   86 │
        │ 2020 │ 107 │   50 │
        └──────┴─────┴──────┘
`,
	},
	{
		style: UnicodeLight,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        ┌──────┬─────┬──────┐
        │ 2018 │ 100 │ 9000 │
        │ 2019 │ 110 │ 85   │
        │      │     │ 86   │
        │      │     │ 86   │
        │ 2020 │ 107 │ 50   │
        └──────┴─────┴──────┘
`,
	},
	{
		style: UnicodeLight,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        ┌──────┬─────┬──────┐
        │ 2018 │ 100 │ 9000 │
        │      │     │  85  │
        │ 2019 │ 110 │  86  │
        │      │     │  86  │
        │ 2020 │ 107 │  50  │
        └──────┴─────┴──────┘
`,
	},
	{
		style: UnicodeLight,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        ┌──────┬─────┬──────┐
        │ 2018 │ 100 │ 9000 │
        │      │     │   85 │
        │      │     │   86 │
        │ 2019 │ 110 │   86 │
        │ 2020 │ 107 │   50 │
        └──────┴─────┴──────┘
`,
	},
	{
		style: UnicodeBold,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        ┏━━━━━━┳━━━━━┳━━━━━━┓
        ┃ 2018 ┃ 100 ┃ 9000 ┃
        ┃ 2019 ┃ 110 ┃ 85   ┃
        ┃      ┃     ┃ 86   ┃
        ┃      ┃     ┃ 86   ┃
        ┃ 2020 ┃ 107 ┃ 50   ┃
        ┗━━━━━━┻━━━━━┻━━━━━━┛
`,
	},
	{
		style: UnicodeBold,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        ┏━━━━━━┳━━━━━┳━━━━━━┓
        ┃ 2018 ┃ 100 ┃ 9000 ┃
        ┃      ┃     ┃  85  ┃
        ┃ 2019 ┃ 110 ┃  86  ┃
        ┃      ┃     ┃  86  ┃
        ┃ 2020 ┃ 107 ┃  50  ┃
        ┗━━━━━━┻━━━━━┻━━━━━━┛
`,
	},
	{
		style: UnicodeBold,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        ┏━━━━━━┳━━━━━┳━━━━━━┓
        ┃ 2018 ┃ 100 ┃ 9000 ┃
        ┃      ┃     ┃   85 ┃
        ┃      ┃     ┃   86 ┃
        ┃ 2019 ┃ 110 ┃   86 ┃
        ┃ 2020 ┃ 107 ┃   50 ┃
        ┗━━━━━━┻━━━━━┻━━━━━━┛
`,
	},
	{
		style: Colon,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        2018 : 100 : 9000
        2019 : 110 : 85
             :     : 86
             :     : 86
        2020 : 107 : 50
`,
	},
	{
		style: Colon,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        2018 : 100 : 9000
             :     :  85
        2019 : 110 :  86
             :     :  86
        2020 : 107 :  50
`,
	},
	{
		style: Colon,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        2018 : 100 : 9000
             :     :   85
             :     :   86
        2019 : 110 :   86
        2020 : 107 :   50
`,
	},
	{
		style: Simple,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
        2019 110 85
                 86
                 86
        2020 107 50
`,
	},
	{
		style: Simple,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
                  85
        2019 110  86
                  86
        2020 107  50
`,
	},
	{
		style: Simple,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
                   85
                   86
        2019 110   86
        2020 107   50
`,
	},
	{
		style: SimpleUnicode,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
        2019 110 85
                 86
                 86
        2020 107 50
`,
	},
	{
		style: SimpleUnicode,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
                  85
        2019 110  86
                  86
        2020 107  50
`,
	},
	{
		style: SimpleUnicode,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
                   85
                   86
        2019 110   86
        2020 107   50
`,
	},
	{
		style: SimpleUnicodeBold,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
        2019 110 85
                 86
                 86
        2020 107 50
`,
	},
	{
		style: SimpleUnicodeBold,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
                  85
        2019 110  86
                  86
        2020 107  50
`,
	},
	{
		style: SimpleUnicodeBold,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        2018 100 9000
                   85
                   86
        2019 110   86
        2020 107   50
`,
	},
	{
		style: Github,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        | 2018 | 100 | 9000 |
        | 2019 | 110 | 85   |
        |      |     | 86   |
        |      |     | 86   |
        | 2020 | 107 | 50   |
`,
	},
	{
		style: Github,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        | 2018 | 100 | 9000 |
        |      |     |  85  |
        | 2019 | 110 |  86  |
        |      |     |  86  |
        | 2020 | 107 |  50  |
`,
	},
	{
		style: Github,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        | 2018 | 100 | 9000 |
        |      |     |   85 |
        |      |     |   86 |
        | 2019 | 110 |   86 |
        | 2020 | 107 |   50 |
`,
	},
	{
		style: CSV,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        2018,100,9000
        2019,110,85
        ,,86
        ,,86
        2020,107,50
`,
	},
	{
		style: CSV,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        2018,100,9000
        ,,85
        2019,110,86
        ,,86
        2020,107,50
`,
	},
	{
		style: CSV,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        2018,100,9000
        ,,85
        ,,86
        2019,110,86
        2020,107,50
`,
	},
	{
		style: JSON,
		align: TL,
		input: borderTestBodyOnly,
		result: `
        {"2018":["100","9000"],"2019":["110","85\n86\n86"],"2020":["107","50"]}
`,
	},
	{
		style: JSON,
		align: MC,
		input: borderTestBodyOnly,
		result: `
        {"2018":["100","9000"],"2019":["110","85\n86\n86"],"2020":["107","50"]}
`,
	},
	{
		style: JSON,
		align: BR,
		input: borderTestBodyOnly,
		result: `
        {"2018":["100","9000"],"2019":["110","85\n86\n86"],"2020":["107","50"]}
`,
	},

	// Empty
	{
		style:  Plain,
		align:  TL,
		input:  ``,
		result: ``,
	},

	// CSV escape
	{
		style: CSV,
		align: TL,
		input: `Year,Income,Source|2018,100,Salary|2019,110,"Consultation"|2020,120,Lottery
et al`,
		rowSep: "|",
		result: `
        Year,Income,Source
        2018,100,Salary
        2019,110,"""Consultation"""
        2020,120,"Lottery
        et al"
`,
	},

	// Missing columns
	{
		style: Unicode,
		align: TL,
		input: `Year,Value
2018,100
2019,
2020,100,200`,
		result: `
        ┏━━━━━━┳━━━━━━━┳━━━━━┓
        ┃ Year ┃ Value ┃     ┃
        ┡━━━━━━╇━━━━━━━╇━━━━━┩
        │ 2018 │ 100   │     │
        │ 2019 │       │     │
        │ 2020 │ 100   │ 200 │
        └──────┴───────┴─────┘
`,
	},
}

func tab(style Style, align Align, data, rowSep string) string {
	tab := New(style)

	rows := strings.Split(data, rowSep)
	if len(rows[0]) > 0 {
		for _, hdr := range strings.Split(rows[0], ",") {
			tab.Header(hdr).SetAlign(align)
		}
	} else {
		var width int
		for _, row := range rows[1:] {
			if len(row) > width {
				width = len(row)
			}
		}
		for i := 0; i < width; i++ {
			tab.SetDefaults(i, align)
		}
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

func m(a, b string) bool {
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

func match(t *testing.T, a, b, name string) {
	if !m(a, b) {
		t.Errorf("%s: got:\n%s\nexpected:\n%s\n", name, a, b)
	}
}

func TestStyles(t *testing.T) {
	for idx, test := range borderTests {
		rowSep := test.rowSep
		if len(rowSep) == 0 {
			rowSep = "\n"
		}
		result := tab(test.style, test.align, test.input, rowSep)
		match(t, result, test.result,
			fmt.Sprintf("TestStyles %d (%s/%s)", idx, test.style, test.align))
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

	var sb strings.Builder
	tab.Print(&sb)
	expected := `
        ┏━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
        ┃     Key ┃            Value             ┃
        ┡━━━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┩
        │    Name │          ACME Corp.          │
        │         │ ┏━━━━━━┳━━━━━━━━┳━━━━━━━━━━┓ │
        │         │ ┃ Year ┃ Income ┃ Expenses ┃ │
        │         │ ┡━━━━━━╇━━━━━━━━╇━━━━━━━━━━┩ │
        │ Numbers │ │ 2018 │    100 │       90 │ │
        │         │ │ 2019 │    110 │       85 │ │
        │         │ │ 2020 │    107 │       50 │ │
        │         │ └──────┴────────┴──────────┘ │
        └─────────┴──────────────────────────────┘
`

	match(t, sb.String(), expected, "TestNested")
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

	var sb strings.Builder
	tab.Print(&sb)
	expected := `
        +----+----+------+
        | 我 | 是 | 测试 |
        +----+----+------+
        | a  | a  | a    |
        +----+----+------+
`

	match(t, sb.String(), expected, "TestWide")
}
