package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func hexStringFromColor(c tcell.Color) string {
	r, g, b := c.RGB()
	return fmt.Sprintf("[#%02x%02x%02x]", r, g, b)
}
