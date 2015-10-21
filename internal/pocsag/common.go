package pocsag

import (
	"github.com/fatih/color"
)

var (
	DEBUG bool
)

var (
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
	blue  = color.New(color.FgBlue)
)

// Tell the package to print debug data
func SetDebug(d bool) {
	DEBUG = d
}
