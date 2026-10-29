package tui

import (
	"github.com/rivo/tview"
)

const (
	// x offset for dialogs that pop up at the selected game
	dialogXOffset = 7

	// dialogs anchored to a row may overlap the keymap footer below the content area
	footerOverlap = 5
)

// modalAt places p as a box of the given size at x/y within the surrounding page
func modalAt(p tview.Primitive, x, y, width, height int) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, y, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, x, 0, false).
			AddItem(p, width, 0, true).
			AddItem(nil, 0, 1, false),
			height, 0, true).
		AddItem(nil, 0, 1, false)
}

// modalCentered places p with the given height vertically centered, using the full width
func modalCentered(p tview.Primitive, height int) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(p, height, 0, true).
		AddItem(nil, 0, 1, false)
}

// anchoredY returns y as position for a box of given height right below an anchor row
// though, if it flows out of the container, the box gets placed on top of the anchor
func anchoredY(y, height int, container *tview.Box) int {
	_, _, _, containerHeight := container.GetRect()
	if y+height > containerHeight+footerOverlap {
		return y - height - 1
	}
	return y
}
