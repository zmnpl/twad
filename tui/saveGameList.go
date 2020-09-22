package tui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/games"
)

const (
	savesHeader = "Savegames"
)

func makeSavegameList(g *games.Game) (*tview.Flex, error) {
	// surrounding container
	frameFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	frameFlex.SetBorderPadding(0, 0, 1, 1)
	frameFlex.AddItem(tview.NewTextView().
		SetText(savesHeader).
		SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)

	// list
	savegameList := tview.NewList()
	frameFlex.AddItem(savegameList, 0, 1, true)
	savegameList.SetSecondaryTextColor(tview.Styles.TitleColor).SetSelectedFocusOnly(true)

	// get savegames
	savegames := g.LoadSavegames()
	fmt.Println(savegames)

	if len(savegames) == 0 {
		return nil, fmt.Errorf("no savegames available")
	}

	// how to populate the list
	populate := func() {
		savegameList.Clear()
		for _, savegame := range savegames {
			savegameList.AddItem("$given_name", fmt.Sprintf("%v (%v)", savegame.FI.Name(), savegame.FI.ModTime().Format("2006-01-02 15:04")), '|', nil)
		}
	}

	// do it
	if savegames != nil {
		populate()
	}

	// hit enter plays demo
	savegameList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// TODO: anything?
	})

	// removes demo at given index and focuses app properly
	removeDemo := func(i int) {
		// TODO: bug in tview; remove when fixed
		savegameList.SetChangedFunc(nil) // BUG WORKAROUND

		//savegames, err = g.RemoveDemo(savegames[i].Name())
		//if err != nil {
		//	showError("could not remove demo", err.Error(), nil, nil)
		//	return
		//}

		if savegames != nil && len(savegames) != 0 {
			populate()
			app.SetFocus(savegameList)
		} else {
			appModeNormal()
		}
		games.Persist()

		//demoList.SetChangedFunc(changeFunc) // BUG WORKAROUND
	}

	// tab navigates back to games table; tab navigation on list is redundant
	savegameList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		// delete mod
		if k == tcell.KeyDelete {
			if savegameList.GetItemCount() > 0 {
				// when in edit mode, this is only confusing
				ci := savegameList.GetCurrentItem()
				if cfg.Instance().DeleteWithoutWarning {
					removeDemo(ci)
					return nil
				}

				youSure := makeYouSureBox(savegames[ci].FI.Name(), // TODO: replace name
					func() {
						removeDemo(ci)
						detailSidePagesSub1.RemovePage(pageYouSure)
					},
					func() {
						detailSidePagesSub1.RemovePage(pageYouSure)
						app.SetFocus(savegameList)
					},
					2, 2, savegameList.Box)
				detailSidePagesSub1.AddPage(pageYouSure,
					youSure, true, true) // TODO: calculate offsets
				app.SetFocus(youSure)
			}
			return nil
		}

		if k == tcell.KeyRune {
			switch event.Rune() {
			// quit app from here as well
			case 'q':
				app.Stop()
				return nil

			// start zip import from here
			case 'i':
				contentPages.SwitchToPage(pageZipImport)
				app.SetFocus(zipInput.selectTree)
				return nil

			}

		}

		return event
	})

	return frameFlex, nil
}
