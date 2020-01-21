package tui

import "github.com/rivo/tview"

const (
	youSureOk = "Yep"
	youSureNo = "Hell No!"
)

// help for navigation
func makeYouSureBox(descriptionText string, onOk func(), xOffset int, yOffset int) *tview.Flex {

	youSureForm := tview.NewForm().
		AddButton(youSureOk, func() {
			onOk()
			appModeNormal()
		}).
		AddButton(youSureNo, appModeNormal)
	youSureForm.SetBorder(false)
	youSureForm.SetFocus(1)

	youSureForm.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)

	description := tview.NewTextView().SetText(descriptionText)
	description.SetBorderPadding(0, 0, 1, 1)
	description.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)

	youSureBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(description, 1, 0, false).
		AddItem(youSureForm, 3, 0, true)
	youSureBox.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)

	//youSureBox.SetBorder(true)

	youSureLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, yOffset, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, xOffset, 1, false).
			AddItem(youSureBox, 50, 0, true).
			AddItem(nil, 0, 1, false),
			4, 1, true).
		AddItem(nil, 0, 1, false)

	return youSureLayout
}
