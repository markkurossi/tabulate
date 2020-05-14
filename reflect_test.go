//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"os"
	"testing"
)

type Outer struct {
	Name    string
	Age     int
	Address *Address
	Info    []Info
}

type Address struct {
	Lines []string
}

type Info struct {
	Email string
	Work  bool
}

func TestReflect(t *testing.T) {
	tab := NewUnicode()
	tab.Header("Field").SetAlign(MR)
	tab.Header("Value").SetAlign(ML)

	err := Reflect(tab, &Outer{
		Name: "Alyssa P. Hacker",
		Age:  45,
		Address: &Address{
			Lines: []string{"42 Hacker way", "03139 Cambridge", "MA"},
		},
		Info: []Info{
			Info{
				Email: "mtr@iki.fi",
			},
			Info{
				Email: "markku.rossi@gmail.com",
				Work:  true,
			},
		},
	})
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}

	tab.Print(os.Stdout)
}
