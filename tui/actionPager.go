package tui

import (
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

// action pager, which holds stats and the "new" form
func makeActionPager() *tview.Pages {
	actionPager = tview.NewPages()
	actionPager.SetTitleAlign(tview.AlignLeft)

	actionPager.AddPage(pageStats, makeStatsTable(&games.Game{}), true, true)

	licensePage = makeLicense()
	actionPager.AddPage(pageLicense, licensePage, true, false)

	return actionPager
}

func makeModListPager() *tview.Pages {
	modListPager := tview.NewPages()
	modListPager.AddPage(pageMods, makeModList(&games.Game{}), true, true)

	return modListPager
}
