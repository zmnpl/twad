package tui

import (
	"path"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/tui/shared"
)

func (u *ui) makeModList(g *games.Game) *tview.Flex {
	modListFlex, modList := newTitledList(shared.Dict.OverviewMods)
	// populate list with data
	for _, mod := range g.Mods {
		colorPrefix := ""
		if !base.ModOk(mod) {
			colorPrefix = "[red]"
		}
		modList.AddItem(colorPrefix+path.Base(mod), path.Dir(mod), 0, nil)
	}

	// edit functionality
	editMode := false
	editOn := func() {
		modList.SetSelectedBackgroundColor(tview.Styles.TertiaryTextColor)
		editMode = true
	}
	editOff := func(save bool) {
		if editMode {
			modList.SetSelectedBackgroundColor(tview.Styles.PrimaryTextColor)
			editMode = false
			if save {
				games.Persist()
			}
		}
	}

	modList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if !editMode {
			editOn()
			return
		}
		editOff(true)
	})

	last := 0
	changeFunc := func(index int, mainText string, secondaryText string, shortcut rune) {
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
	}
	modList.SetChangedFunc(changeFunc)

	removeMod := func(i int) {
		// TODO: bug in tview
		// Existing change func when deleting zero item
		// created pull request; setting nil and resetting is temp workaround
		modList.SetChangedFunc(nil) // BUG WORKAROUND
		modList.RemoveItem(i)
		g.RemoveMod(i)
		games.Persist()
		modList.SetChangedFunc(changeFunc) // BUG WORKAROUND
	}

	// tab navigates back to games table; tab navigation on list is redundant
	modList.SetInputCapture(u.withQuit(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		// switch back to game table
		if k == tcell.KeyTab {
			u.app.SetFocus(u.gamesTable)
			return nil
		}

		// delete mod
		if k == tcell.KeyDelete {
			// when in edit mode, this is only confusing
			if !editMode && modList.GetItemCount() > 0 {
				ci := modList.GetCurrentItem()
				if base.Config().DeleteWithoutWarning {
					removeMod(ci)
					return nil
				}

				confirm(u.detailSidePagesSub1, g.Mods[ci], 2, 2, modList, func() { // TODO: calculate offsets
					removeMod(ci)
				})
			}
			return nil
		}

		// add mod
		if k == tcell.KeyRune && event.Rune() == 'm' {
			u.showModTree(g)
			return nil
		}

		return event
	}))

	return modListFlex
}
