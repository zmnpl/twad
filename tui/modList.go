package tui

import (
	"path"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	overviewMods = "Mods in order"
)

func populateModList(list *tview.List, g *games.Game) {
	i := 0
	for _, mod := range g.Mods {
		i++
		list.AddItem(path.Base(mod), path.Dir(mod), '*', nil)
	}
}

func makeModList(g *games.Game) *tview.Flex {

	modListFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	modListFlex.AddItem(tview.NewTextView().
		SetText(overviewMods).
		SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)

	modList := tview.NewList()
	modList.SetSecondaryTextColor(tview.Styles.TitleColor).SetSelectedFocusOnly(true)
	populateModList(modList, g)

	//mover := func(selectedGame *games.Game) func() *tview.TreeView {
	//	return func() *tview.TreeView {
	//		return makeModTree(selectedGame)
	//	}
	//}

	editMode := false
	modList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if editMode == false {
			modList.SetSelectedBackgroundColor(tview.Styles.TertiaryTextColor)
			editMode = true
			return
		}
		modList.SetSelectedBackgroundColor(tview.Styles.PrimaryTextColor)
		editMode = false
	})

	modList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {

	})

	// tab navigates back to games table; tab navigation on list is redundant
	modList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyTab {
			app.SetFocus(gamesTable)
			return nil
		}

		return event
	})

	modListFlex.AddItem(modList, 0, 1, true)

	return modListFlex
}
