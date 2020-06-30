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
		tab.Header(hdr).SetAlign(align)
	}

	for i := 1; i < len(rows); i++ {
		row := tab.Row()
		for _, col := range strings.Split(rows[i], ",") {
			row.ColumnData(NewText(col))
		}
	}
	return tab
}

func tabulate(tab *Tabulate, align Align) *Tabulate {
	return tabulateRows(tab, align, strings.Split(data, "\n"))
}

func align(align Align) {
	tabulate(New(Plain), align).Print(os.Stdout)
	tabulate(New(ASCII), align).Print(os.Stdout)
	tabulate(New(Unicode), align).Print(os.Stdout)
	tabulate(New(Colon), align).Print(os.Stdout)
	tabulate(New(Simple), align).Print(os.Stdout)
	tabulate(New(Github), align).Print(os.Stdout)
	tabulate(New(JSON), align).Print(os.Stdout)
}

func TestBorders(t *testing.T) {
	align(TL)
	align(MC)
	align(BR)
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
	row.ColumnData(tabulate(New(Unicode), TR))

	tab.Print(os.Stdout)
}
