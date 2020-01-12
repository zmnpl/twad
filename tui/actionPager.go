package tui

import "github.com/rivo/tview"

// action pager, which holds stats and the "new" form
func makeActionPager() *tview.Pages {
	actionPager = tview.NewPages()
	actionPager.SetTitleAlign(tview.AlignLeft)

	statsTable = makeStatsTable()
	actionPager.AddPage(pageStats, statsTable, true, true)

	newForm = makeNewGameForm()
	actionPager.AddPage(pageNewForm, newForm, true, false)

	licensePage = makeLicense()
	actionPager.AddPage(pageLicense, licensePage, true, false)

	return actionPager
}