package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/tui/shared"
)

// command preview
func (u *ui) makeCommandPreview() *tview.TextView {
	u.commandPreview = tview.NewTextView().SetDynamicColors(true)
	u.commandPreview.SetBorder(true)

	return u.commandPreview
}

func (u *ui) populateCommandPreview(g *games.Game) {
	u.commandPreview.Clear()
	fmt.Fprintf(u.commandPreview, "%s %s", shared.PreviewText, shared.StylizeCommandList(g.CommandList(), colorTagMoreContrast, colorTagContrast, colorTagPrimaryText))
}
