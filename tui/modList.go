package tui

import (
	"path"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	overviewMods = "Mods"
)

func makeModList(g *games.Game) *tview.Flex {

	modListFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	modListFlex.AddItem(tview.NewTextView().SetText(overviewMods).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)
	modList := tview.NewList()
	i := 0
	for _, mod := range g.Mods {
		i++
		modList.AddItem(path.Base(mod), path.Dir(mod), '*', nil)
	}
	modListFlex.AddItem(modList, 0, 1, false)

	return modListFlex
}
