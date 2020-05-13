//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

type Format int

const (
	FmtNone Format = iota
	FmtBold
	FmtItalic
)

func (fmt Format) VT100() string {
	switch fmt {
	case FmtBold:
		return "\x1b[1m"
	case FmtItalic:
		return "\x1b[3m"
	default:
		return "\x1b[m"
	}
}
