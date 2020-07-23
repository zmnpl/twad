package tui

import "github.com/rivo/tview"

const (
	confirmText = "Yep"
	abortText   = "Hell No"
)

// help for navigation
func makeYouSureBox(question string, onOk func(), xOffset int, yOffset int) *tview.Flex {

	youSureForm := tview.NewForm().
		AddButton(confirmText, func() {
			onOk()
			appModeNormal()
		}).
		AddButton(abortText, appModeNormal)
	youSureForm.SetBorder(false)
	youSureForm.SetFocus(1)

	//	youSureForm.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)

	description := tview.NewTextView().SetText(question)
	description.SetBorderPadding(1, 0, 1, 1)
	//	description.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)

	youSureBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(description, 2, 0, false).
		AddItem(youSureForm, 3, 0, true)
		//	youSureBox.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)
	youSureBox.SetBorder(true)

	width := len(question) + 4
	minWidth := 11 + len(confirmText) + len(abortText)
	if width < minWidth {
		width = minWidth
	}

	// TODO: catch yOffset if popup flows out of the windows

	youSureLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, yOffset, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, xOffset, 1, false).
			AddItem(youSureBox, width, 0, true).
			AddItem(nil, 0, 1, false),
			7, 1, true).
		AddItem(nil, 0, 1, false)

	return youSureLayout
}
