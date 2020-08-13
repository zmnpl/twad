package tui

import (
	"path"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/games"
)

const (
	overviewMods = "Mods in order"
)

func makeModList(g *games.Game) *tview.Flex {
	// surrounding container
	modListFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	modListFlex.SetBorderPadding(0, 0, 1, 1)
	modListFlex.AddItem(tview.NewTextView().
		SetText(overviewMods).
		SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)

	// list
	modList := tview.NewList()
	modListFlex.AddItem(modList, 0, 1, true)
	modList.SetSecondaryTextColor(tview.Styles.TitleColor).SetSelectedFocusOnly(true)
	// populate list with data
	for _, mod := range g.Mods {
		modList.AddItem(path.Base(mod), path.Dir(mod), '|', nil)
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
		if editMode == false {
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
	modList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		// switch back to game table
		if k == tcell.KeyTab {
			app.SetFocus(gamesTable)
			return nil
		}

		// delete mod
		if k == tcell.KeyDelete {
			// when in edit mode, this is only confusing
			if !editMode && modList.GetItemCount() > 0 {
				ci := modList.GetCurrentItem()
				if cfg.GetInstance().DeleteWithoutWarning {
					removeMod(ci)
					return nil
				}

				detailSidePagesSub1.AddPage(pageYouSure,
					makeYouSureBox(*g,
						func() {
							removeMod(ci)
							detailSidePagesSub1.RemovePage(pageYouSure)
							app.SetFocus(modList)
						},
						func() {
							//appModeNormal()
							detailSidePagesSub1.RemovePage(pageYouSure)
							app.SetFocus(modList)
						},
						2, 2, modList.Box), true, true) // TODO: calculate offsets
			}
			return nil
		}

		if k == tcell.KeyRune {
			switch event.Rune() {
			// add mod
			case 'm':
				modTree := makeModTree(g)
				detailSidePagesSub2.AddPage(pageModSelector, modTree, true, false)
				detailSidePagesSub2.SwitchToPage(pageModSelector)
				app.SetFocus(modTree)
				return nil

			// quit app from here as well
			case 'q':
				app.Stop()
				return nil
			}
		}

		return event
	})

	return modListFlex
}
