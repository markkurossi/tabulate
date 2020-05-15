//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"fmt"
	"os"
	"testing"
)

type Outer struct {
	Name    string
	Comment string `tabulate:"@detail"`
	Age     int
	Address *Address `tabulate:"omitempty"`
	Info    []*Info
}

type Address struct {
	Lines []string
}

type Info struct {
	Email string
	Work  bool
}

func reflectTest(flags Flags, tags []string, v interface{}) error {
	tab := NewUnicode()
	tab.Header("Field").SetAlign(MR)
	tab.Header("Value").SetAlign(ML)

	err := Reflect(tab, flags, tags, v)
	if err != nil {
		return err
	}

	tab.Print(os.Stdout)
	return nil
}

func TestReflect(t *testing.T) {
	err := reflectTest(OmitEmpty, nil, &Outer{
		Name: "Alyssa P. Hacker",
		Age:  45,
		Address: &Address{
			Lines: []string{"42 Hacker way", "03139 Cambridge", "MA"},
		},
		Info: []*Info{
			&Info{
				Email: "mtr@iki.fi",
			},
			&Info{
				Email: "markku.rossi@gmail.com",
				Work:  true,
			},
		},
	})
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}

	data := &Outer{
		Name:    "Alyssa P. Hacker",
		Comment: "Structure and Interpretation of Computer Programs",
		Age:     45,
		Info: []*Info{
			nil,
			&Info{
				Email: "markku.rossi@gmail.com",
				Work:  true,
			},
		},
	}

	err = reflectTest(OmitEmpty, nil, data)
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
	err = reflectTest(0, []string{"detail"}, data)
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
}

type Outer2 struct {
	Name  string
	Inner *Inner
}

type Inner struct {
	A int
	B int
}

func (in Inner) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("A=%v, B=%v", in.A, in.B)), nil
}

func TestReflectTextMarshaler(t *testing.T) {
	err := reflectTest(0, nil, &Outer2{
		Name: "ACME Corp.",
		Inner: &Inner{
			A: 100,
			B: 42,
		},
	})
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
}
