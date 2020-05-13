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

type Outer struct {
	Name    string
	Age     int
	Address *Address
	Info    Info
}

type Address struct {
	Street string
	Zip    string
}

type Info struct {
	Email     string
	Permanent bool
}

func TestReflect(t *testing.T) {
	tab := NewUnicode()
	tab.Header(Right, Middle, NewText("Field"))
	tab.Header(Left, Middle, NewText("Value"))

	err := Reflect(tab, &Outer{
		Name: "Alyssa P. Hacker",
		Age:  45,
		Address: &Address{
			Street: "Hacker way",
			Zip:    "02139",
		},
		Info: Info{
			Email: "mtr@iki.fi",
		},
	})
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}

	tab.Print(os.Stdout)
}

var data = `Year,Income,Expenses
2018,100,90
2019,110,85
2020,107,50`

func tabulateRows(tab *Tabulate, align Align, valign VAlign,
	rows []string) *Tabulate {

	for _, hdr := range strings.Split(rows[0], ",") {
		tab.Header(align, valign, NewText(hdr))
	}

	for i := 1; i < len(rows); i++ {
		row := tab.Row()
		for _, col := range strings.Split(rows[i], ",") {
			row.Column(NewText(col))
		}
	}
	return tab
}

func tabulate(tab *Tabulate, align Align, valign VAlign) *Tabulate {
	return tabulateRows(tab, align, valign, strings.Split(data, "\n"))
}

func align(align Align, valign VAlign) {
	tabulate(NewWS(), align, valign).Print(os.Stdout)
	tabulate(NewASCII(), align, valign).Print(os.Stdout)
	tabulate(NewUnicode(), align, valign).Print(os.Stdout)
	tabulate(NewColon(), align, valign).Print(os.Stdout)
}

func TestBorders(t *testing.T) {
	align(Left, Top)
	align(Center, Middle)
	align(Right, Bottom)
}

var csv = `Year,Income,Source|2018,100,Salary|2019,110,"Consultation"|2020,120,Lottery
et al`

func TestCSV(t *testing.T) {
	rows := strings.Split(csv, "|")
	tabulateRows(NewCSV(), None, Top, rows).Print(os.Stdout)
}

func TestNested(t *testing.T) {
	tab := NewUnicode()

	tab.Header(Right, Middle, NewLines("Key"))
	tab.Header(Center, Middle, NewLines("Value"))

	row := tab.Row()
	row.Column(NewLines("Name"))
	row.Column(NewLines("ACME Corp."))

	row = tab.Row()
	row.Column(NewLines("Numbers"))
	row.Column(tabulate(NewUnicode(), Right, Top).Data())

	tab.Print(os.Stdout)
}
