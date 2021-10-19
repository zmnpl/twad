package tui

import (
	"github.com/rivo/tview"
)

// help for navigation
func makeYouSureBox(title string, onOk func(), onCancel func(), xOffset int, yOffset int, container *tview.Box) *tview.Flex {
	youSureForm := tview.NewForm().
		AddButton(dict.confirmText, onOk).
		AddButton(dict.abortText, onCancel)
	youSureForm.
		SetBorder(true).
		SetTitle(title)
	youSureForm.SetFocus(1)

	height := 5
	width := 50

	// surrounding layout
	_, _, _, containerHeight := container.GetRect()
	helpHeight := 5

	// default: right below the selected game
	// though, if it flows out of the screen, then on top of the game
	if yOffset+height > containerHeight+helpHeight {
		yOffset = yOffset - height - 1
	}

	youSureLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, yOffset, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, xOffset, 0, false).
			AddItem(youSureForm, width, 0, true).
			AddItem(nil, 0, 1, false),
			height, 1, true).
		AddItem(nil, 0, 1, false)

	return youSureLayout
}
