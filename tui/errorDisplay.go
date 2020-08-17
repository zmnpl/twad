package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	errConfirm = "I don't care. Go ahead."
	errAbort   = "Let me fix that first!"
)

// help for navigation
func makeErrorDisplay(errTitle string, errString string, response chan bool) *tview.Flex {
	errForm := tview.NewForm().
		AddButton(errConfirm, func() {
			response <- true
		}).
		AddButton(errAbort, func() {
			response <- false
		})
	errForm.
		SetBorder(true).
		SetTitle(errTitle).SetBackgroundColor(tcell.ColorRed)
	errForm.SetFocus(1)

	height := 10
	errLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(errForm, height, 0, true).
		AddItem(nil, 0, 1, false)

	return errLayout
}
