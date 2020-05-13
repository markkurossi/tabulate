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

func tabulateRows(tab *Tabulate, align Align, rows []string) *Tabulate {
	for _, hdr := range strings.Split(rows[0], ",") {
		tab.Header(align, NewText(hdr))
	}

	for i := 1; i < len(rows); i++ {
		row := tab.Row()
		for _, col := range strings.Split(rows[i], ",") {
			row.Column(NewText(col))
		}
	}
	return tab
}

func tabulate(tab *Tabulate, align Align) *Tabulate {
	return tabulateRows(tab, align, strings.Split(data, "\n"))
}

func align(align Align) {
	tabulate(NewTabulateWS(), align).Print(os.Stdout)
	tabulate(NewTabulateASCII(), align).Print(os.Stdout)
	tabulate(NewTabulateUnicode(), align).Print(os.Stdout)
	tabulate(NewTabulateColon(), align).Print(os.Stdout)
}

func TestBorders(t *testing.T) {
	align(AlignLeft)
	align(AlignCenter)
	align(AlignRight)
}

var csv = `Year,Income,Source|2018,100,Salary|2019,110,"Consultation"|2020,120,Lottery
et al`

func TestCSV(t *testing.T) {
	rows := strings.Split(csv, "|")
	tabulateRows(NewTabulateCSV(), AlignNone, rows).Print(os.Stdout)
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
