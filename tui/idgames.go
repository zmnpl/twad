package tui

import (
	"github.com/rivo/tview"
	"github.com/zmnpl/goidgames"
)

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

func makeIdgamesBrowser() *IdgamesBrowser {
	layout := tview.NewGrid()
	layout.SetRows(1, -1)
	layout.SetColumns(-1, -1)

	list := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')
	gamesTable.SetBorderPadding(0, 0, 1, 2)

	list.SetCell(0, 0, tview.NewTableCell(gameTableHeaderRating).SetTextColor(tview.Styles.SecondaryTextColor))
	list.SetCell(0, 1, tview.NewTableCell(gameTableHeaderName).SetTextColor(tview.Styles.SecondaryTextColor))
	list.SetCell(0, 2, tview.NewTableCell(gameTableHeaderSourcePort).SetTextColor(tview.Styles.SecondaryTextColor))
	list.SetCell(0, 3, tview.NewTableCell(gameTableHeaderIwad).SetTextColor(tview.Styles.SecondaryTextColor))

	idgameFiles, _ := goidgames.LatestFiles(10, 0)
	fixRows, fixCols := 1, 4
	rows, cols := len(idgameFiles), 0

	for r := 1; r < rows+fixRows; r++ {
		var f goidgames.Idgame
		if r > 0 {
			f = idgameFiles[r-fixRows]
		}
		for c := 0; c < cols+fixCols; c++ {
			var cell *tview.TableCell

			switch c {
			case 0:
				cell = tview.NewTableCell(f.Title).SetTextColor(tview.Styles.TitleColor)
			case 1:
				cell = tview.NewTableCell("test").SetTextColor(tview.Styles.SecondaryTextColor)
			case 2:
				cell = tview.NewTableCell("test").SetTextColor(tview.Styles.PrimaryTextColor)
			case 3:
				cell = tview.NewTableCell("test").SetTextColor(tview.Styles.PrimaryTextColor)
			default:
				cell = tview.NewTableCell("").SetTextColor(tview.Styles.PrimaryTextColor)
			}

			list.SetCell(r, c, cell)
		}
	}

	layout.AddItem(list, 1, 0, 1, 2, 0, 0, true)

	return &IdgamesBrowser{
		layout: layout,
		list:   list,
	}
}
