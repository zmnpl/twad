package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	errTitleStart = "ERROR"
	errYolo       = "I don't care. Go ahead."
	errNotYet     = "Let me fix that first!"
	errAbort      = "Ok"
)

// can show errors prominently on the screen
// if a function is supplied, the user gets the choice to proceed
// without fixing what might have cause this on his/her own risk
func showError(errTitle string, errString string, handFocusBackTo tview.Primitive, YOLO func()) {
	// form with buttons
	errForm := tview.NewForm()

	// sets focus to the given primitive
	// if nil was given, then the apps default state will be restored
	resetFocus := func() {
		if handFocusBackTo != nil {
			app.SetFocus(handFocusBackTo)
		} else {
			appModeNormal()
		}
	}

	// YOLO lets the user execute and action despite the error
	if YOLO != nil {
		errForm.AddButton(errYolo, func() {
			YOLO()
			contentPages.RemovePage(pageError)
			resetFocus()
		})

		errForm.AddButton(errNotYet, func() {
			contentPages.RemovePage(pageError)
			resetFocus()
		})
	} else {
		errForm.AddButton(errAbort, func() {
			contentPages.RemovePage(pageError)
			resetFocus()
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
