package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

const (
	previewText = "preview"
)

// command preview
func makeCommandPreview() *tview.TextView {
	commandPreview = tview.NewTextView().
		SetDynamicColors(true)
	commandPreview.SetBackgroundColor(previewBackgroundColor)
	fmt.Fprintf(commandPreview, "")

	return commandPreview
}

func populateCommandPreview(command string) {
	commandPreview.Clear()
	fmt.Fprintf(commandPreview, previewText+" $ %s", command)
}
