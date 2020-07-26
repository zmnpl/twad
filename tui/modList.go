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

func makeModList(g *games.Game) *tview.Flex {

	modListFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	modListFlex.AddItem(tview.NewTextView().
		SetText(overviewMods).
		SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)

	modList := tview.NewList()
	modList.SetSecondaryTextColor(tview.Styles.TitleColor).SetSelectedFocusOnly(true)
	i := 0
	for _, mod := range g.Mods {
		i++
		modList.AddItem(path.Base(mod), path.Dir(mod), '*', nil)
	}
	//mover := func(selectedGame *games.Game) func() *tview.TreeView {
	//	return func() *tview.TreeView {
	//		return makeModTree(selectedGame)
	//	}
	//}

	editMode := false
	editOn := func() {
		modList.SetSelectedBackgroundColor(tview.Styles.TertiaryTextColor)
		editMode = true
	}
	editOff := func() {
		modList.SetSelectedBackgroundColor(tview.Styles.PrimaryTextColor)
		editMode = false
		games.Persist()
	}

	modList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if editMode == false {
			editOn()
			return
		}
		editOff()
	})

	last := 0
	modList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if editMode {
			g.SwitchMods(last, index)
			lastMain, lastSecondary := modList.GetItemText(last)
			main, secondary := modList.GetItemText(index)

			modList.SetItemText(index, lastMain, lastSecondary)
			modList.SetItemText(last, main, secondary)
		}
		last = index
	})

	modList.SetDoneFunc(func() {
		editOff()
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
