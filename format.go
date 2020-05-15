//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package tabulate

// Format specifies text formatting.
type Format int

// Format values specify various VT100 text formatting codes.
const (
	FmtNone Format = iota
	FmtBold
	FmtItalic
)

// VT100 creates VT100 terminal emulation codes for the agument
// format.
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
