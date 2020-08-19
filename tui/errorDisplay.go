package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	errTitleStart = "A pinky demon bit your ass!"
	errYolo       = "I don't care. Go ahead."
	errNotYet     = "Let me fix that first!"
	errAbort      = "Ok"
)

// can show errors prominently on the screen
// if a function is supplied, the user gets the choice to proceed
// without fixing what might have cause this on his/her own risk
func showError(errTitle string, errString string, onIDontCare func()) {
	// form with buttons
	errForm := tview.NewForm()

	if onIDontCare != nil {
		errForm.AddButton(errYolo, func() {
			onIDontCare()
			contentPages.RemovePage(pageError)
		})

		errForm.AddButton(errNotYet, func() {
			contentPages.RemovePage(pageError)
		})
	} else {
		errForm.AddButton(errAbort, func() {
			contentPages.RemovePage(pageError)
		})
	}

	// style
	errForm.SetButtonBackgroundColor(tcell.ColorRed)
	errForm.SetButtonTextColor(tcell.ColorWhite)

	errorText := tview.NewTextView()
	errorText.SetDynamicColors(true)
	errorText.SetText(errString)
	errorText.SetBorderPadding(0, 0, 1, 1)

	foo := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(errorText, 5, 0, false).
		AddItem(errForm, 3, 0, true)
	foo.
		SetBorder(true).
		SetTitle(errTitleStart + " - " + errTitle).
		SetBorderColor(tcell.ColorRed).
		SetTitleColor(tcell.ColorRed)

	errForm.SetFocus(1)
	height := 10
	errLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(foo, height, 0, true).
		AddItem(nil, 0, 1, false)

	contentPages.AddPage(pageError, errLayout, true, true)
}
