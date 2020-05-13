//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"os"
	"strings"
	"testing"
)

var data = `Year,Income,Expenses
2018,100,90
2019,110,85
2020,107,50`

func tabulate(tab *Tabulate, align Align) *Tabulate {
	rows := strings.Split(data, "\n")
	for _, hdr := range strings.Split(rows[0], ",") {
		tab.Header(align, NewLines(hdr))
	}

	for i := 1; i < len(rows); i++ {
		row := tab.Row()
		for _, col := range strings.Split(rows[i], ",") {
			row.Column(NewLines(col))
		}
	}
	return tab
}

func align(align Align) {
	tabulate(NewTabulateWS(), align).Print(os.Stdout)
	tabulate(NewTabulateASCII(), align).Print(os.Stdout)
	tabulate(NewTabulateUnicode(), align).Print(os.Stdout)
}

func TestBorders(t *testing.T) {
	align(AlignLeft)
	align(AlignCenter)
	align(AlignRight)
}

func TestNested(t *testing.T) {
	tab := NewTabulateUnicode()

	tab.Header(AlignRight, NewLines("Key"))
	tab.Header(AlignCenter, NewLines("Value"))

	row := tab.Row()
	row.Column(NewLines("Name"))
	row.Column(NewLines("ACME Corp."))

	row = tab.Row()
	row.Column(NewLines("Numbers"))
	row.Column(tabulate(NewTabulateUnicode(), AlignRight).Data())

	tab.Print(os.Stdout)
}
