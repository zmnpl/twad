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
	width := 64
	youSureLayout := modalAt(youSureForm, xOffset, anchoredY(yOffset, height, container), width, height)

	app.SetFocus(youSureForm)

	return youSureLayout
}
