package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
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
	savegames := g.LoadSavegames() //time.Sleep(2 * time.Second)

	if len(savegames) == 0 || savegames == nil {
		return nil, fmt.Errorf("no savegames available")
	}

	var statsTable *tview.Table

	// how to populate the list
	populate := func() {
		savegameList.Clear()
		for i, savegame := range savegames {
			savegameList.AddItem(savegame.Meta.Title, fmt.Sprintf("%v (%v)", savegame.FI.Name(), savegame.FI.ModTime().Format("2006-01-02 15:04:05")), '|', nil)

			if i == 0 {
				statsTable = makeLevelStatsTable(*savegames[i], savegameList)
				detailSidePagesSub2.AddPage(pageStats, statsTable, true, true)
			}
		}
	}

	populate()

	// load savegame on enter
	savegameList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// TODO: Load savegame on enter
	})

	savegameList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		statsTable = makeLevelStatsTable(*savegames[index], savegameList)
		detailSidePagesSub2.AddPage(pageStats, statsTable, true, true)
	})

	// tab navigates back to games table; tab navigation on list is redundant
	savegameList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		if k == tcell.KeyTAB && statsTable != nil {
			app.SetFocus(detailSidePagesSub2)
		}

		if k == tcell.KeyRune {
			switch event.Rune() {
			// quit app from here as well
			case 'q':
				app.Stop()
				return nil
			}
		}

		return event
	})

	return frameFlex, nil
}
