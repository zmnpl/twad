package tui

import "github.com/rivo/tview"

//"github.com/zmnpl/goidgames"

type IdgamesBrowser struct {
	layout      *tview.Grid
	list        *tview.Table
	fileDetails *tview.Table
	reviews     *tview.Table
	search      *tview.InputField
}

func init() {
	//goidgames.Get(1338, "")
}
