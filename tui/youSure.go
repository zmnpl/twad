package tui

import "github.com/rivo/tview"

const (
	youSureOk = "Yep"
	youSureNo = "Hell No!"
)

// help for navigation
func makeYouSureBox(descriptionText string, onOk func()) *tview.Flex {

	youSureForm := tview.NewForm().
		AddButton(youSureOk, onOk).
		AddButton(youSureNo, appModeNormal)
	youSureForm.SetBorder(false)
	youSureForm.SetFocus(1)

	description := tview.NewTextView().SetText(descriptionText)
	description.SetBorderPadding(1, 1, 1, 1)

	youSureBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(description, 3, 0, false).
		AddItem(youSureForm, 6, 0, true)

	youSureBox.SetBorder(true)

	youSureLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(youSureBox, 25, 0, true).
			AddItem(nil, 0, 1, false),
			0, 1, true).
		AddItem(nil, 0, 1, false)

	return youSureLayout
}
