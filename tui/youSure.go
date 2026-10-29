package tui

import (
	"github.com/rivo/tview"
)

// confirm asks the user to confirm an action with a dialog at x, right below row y of the given layer
// either choice closes the dialog and hands focus back to returnTo; nil restores the apps default state
func confirm(l *layer, title string, x, y int, returnTo tview.Primitive, onOk func()) {
	youSureForm := tview.NewForm().
		AddButton(dict.confirmText, func() {
			l.close(pageYouSure, returnTo)
			onOk()
		}).
		AddButton(dict.abortText, func() {
			l.close(pageYouSure, returnTo)
		})
	youSureForm.
		SetBorder(true).
		SetTitle(title)
	youSureForm.SetFocus(1)

	l.showAt(pageYouSure, youSureForm, x, y, 64, 5)
}
