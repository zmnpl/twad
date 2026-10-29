package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// can show errors prominently on the screen
// if a function is supplied, the user gets the choice to proceed
// without fixing what might have cause this on his/her own risk
func (u *ui) showError(errTitle string, errString string, handFocusBackTo tview.Primitive, YOLO func()) {
	// form with buttons
	errForm := tview.NewForm()

	// sets focus to the given primitive
	// if nil was given, then the apps default state will be restored
	resetFocus := func() {
		u.contentPages.close(pageError, handFocusBackTo)
	}

	// YOLO lets the user execute and action despite the error
	if YOLO != nil {
		errForm.AddButton(dict.errYolo, func() {
			YOLO()
			resetFocus()
		})
		errForm.AddButton(dict.errNotYet, func() {
			resetFocus()
		})
	} else {
		errForm.AddButton(dict.errAbort, func() {
			resetFocus()
		})
	}

	// style
	errForm.SetButtonBackgroundColor(warnColor)
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
		SetTitle(dict.errTitleStart + " - " + errTitle).
		SetBorderColor(warnColor).
		SetTitleColor(warnColor)

	errForm.SetFocus(0)
	u.contentPages.show(pageError, modalCentered(foo, 10))

	// very dirty hack to retrieve focus...
	// TODO: how to do better when produced by form item done funcs
	//getFocus := func() {
	//	time.Sleep(1 * time.Second)
	//	gf := func() {
	//		app.SetFocus(errForm)
	//	}
	//	app.QueueUpdateDraw(gf)
	//}
	//go getFocus()
}
