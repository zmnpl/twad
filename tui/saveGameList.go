package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

func (u *ui) makeSavegameList(g *games.Game) (*tview.Flex, error) {
	frameFlex, savegameList := newTitledList(dict.savesHeader)

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
				statsTable = u.makeLevelStatsTable(*savegames[i], savegameList)
				u.detailSidePagesSub2.add(pageSaveStats, statsTable)
			}
		}
	}

	populate()

	// load savegame on enter
	savegameList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// TODO: Load savegame on enter
	})

	savegameList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		statsTable = u.makeLevelStatsTable(*savegames[index], savegameList)
		u.detailSidePagesSub2.add(pageSaveStats, statsTable)
	})

	// tab navigates back to games table; tab navigation on list is redundant
	savegameList.SetInputCapture(u.withQuit(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB && statsTable != nil {
			u.app.SetFocus(u.detailSidePagesSub2)
		}
		return event
	}))

	return frameFlex, nil
}
