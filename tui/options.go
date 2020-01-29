package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	optionsPathLabel       = "Base Path:"
	optionsWarnBeforeLabel = "Warn before deltion?"
)

func makeOptions() *tview.Flex {
	pathLabel := tview.NewTextView().SetText(optionsPathLabel)
	path := tview.NewInputField()
	pathRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pathLabel, 0, 2, false).
		AddItem(path, 0, 2, true)

	warnBeforeDeleteLabel := tview.NewTextView().SetText(optionsWarnBeforeLabel)
	warn := tview.NewCheckbox()
	warnRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(warnBeforeDeleteLabel, 0, 2, false).
		AddItem(warn, 0, 3, false)

	path.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyTab {
			// TODO
			app.SetFocus(warn)
			return nil
		}

		return event
	})

	options := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pathRow, 0, 1, true).
		AddItem(warnRow, 0, 1, false)
	return options
}
