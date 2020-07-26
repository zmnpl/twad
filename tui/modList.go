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
	// surrounding container
	modListFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	modListFlex.AddItem(tview.NewTextView().
		SetText(overviewMods).
		SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)

	// list
	modList := tview.NewList()
	modList.SetSecondaryTextColor(tview.Styles.TitleColor).SetSelectedFocusOnly(true)
	// populate list with data
	i := 0
	for _, mod := range g.Mods {
		i++
		modList.AddItem(path.Base(mod), path.Dir(mod), '*', nil)
	}

	// edit functionality
	editMode := false
	editOn := func() {
		modList.SetSelectedBackgroundColor(tview.Styles.TertiaryTextColor)
		editMode = true
	}
	editOff := func() {
		if editMode {
			modList.SetSelectedBackgroundColor(tview.Styles.PrimaryTextColor)
			editMode = false
			games.Persist()
		}
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
			// switch mod positions in game
			g.SwitchMods(last, index)

			// switch list item texts
			lastMain, lastSecondary := modList.GetItemText(last)
			main, secondary := modList.GetItemText(index)
			modList.SetItemText(index, lastMain, lastSecondary)
			modList.SetItemText(last, main, secondary)
		}
		last = index
	})

	modList.SetDoneFunc(func() {
		//editOff()
	})

	// tab navigates back to games table; tab navigation on list is redundant
	modList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyTab {
			app.SetFocus(gamesTable)
			return nil
		}

		if k == tcell.KeyDEL {
			modList.RemoveItem(modList.GetCurrentItem())

			editOff()
			modList.RemoveItem(modList.GetCurrentItem())
			// TODO: actually remove mod from the game
			// need to write function on game for that
			return nil
		}

		if k == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				editOff()
				app.Stop()
			}
		}

		return event
	})

	modListFlex.AddItem(modList, 0, 1, true)

	return modListFlex
}
