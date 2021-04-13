package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/goidgames"
)

//"github.com/zmnpl/goidgames"

type IdgamesBrowser struct {
	layout          *tview.Grid
	list            *tview.Table
	fileDetails     *tview.Table
	fileDetailsText *tview.TextView
	reviews         *tview.Table
	search          *tview.InputField
}

func init() {
	//goidgames.Get(1338, "")
}

func makeIdgamesBrowser() *IdgamesBrowser {
	layout := tview.NewGrid()
	layout.SetRows(1, -1)
	layout.SetColumns(-1, -1)

	// list with results
	list := tview.NewTable().
		SetFixed(1, 2).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')
	gamesTable.SetBorderPadding(0, 0, 1, 2)
	layout.AddItem(list, 1, 0, 1, 1, 0, 0, true)

	// details for selection
	details := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	details.SetBorderPadding(0, 0, 1, 1)
	layout.AddItem(details, 1, 1, 1, 1, 0, 0, false)

	return &IdgamesBrowser{
		layout:          layout,
		list:            list,
		fileDetailsText: details,
	}
}

func (browser *IdgamesBrowser) UpdateLatest() {
	go func() {
		app.QueueUpdateDraw(func() {
			idgames, _ := goidgames.LatestFiles(10, 0)

			go func() {
				for i, _ := range idgames {
					g, err := goidgames.Get(idgames[i].Id, "")
					if err != nil {
						continue
					}
					idgames[i] = g
				}
			}()

			browser.populateList(idgames)
		})
	}()
}

func (browser *IdgamesBrowser) populateList(idgames []goidgames.Idgame) {
	browser.list.Clear()
	// header
	browser.list.SetCell(0, 0, tview.NewTableCell("Rating").SetTextColor(tview.Styles.SecondaryTextColor))
	browser.list.SetCell(0, 1, tview.NewTableCell("Title").SetTextColor(tview.Styles.SecondaryTextColor))
	browser.list.SetCell(0, 2, tview.NewTableCell("Author").SetTextColor(tview.Styles.SecondaryTextColor))
	browser.list.SetCell(0, 3, tview.NewTableCell("Date").SetTextColor(tview.Styles.SecondaryTextColor))

	browser.list.SetSelectionChangedFunc(func(r int, c int) {
		switch r {
		case 0:
			return
		default:
			browser.populateDetails(idgames[r-1])
		}
	})

	fixRows := 1
	cols := 4
	rows := len(idgames)
	for r := 1; r < rows+fixRows; r++ {
		var f goidgames.Idgame
		if r > 0 {
			f = idgames[r-fixRows]
		}
		for c := 0; c < cols; c++ {
			var cell *tview.TableCell

			switch c {
			case 0:
				cell = tview.NewTableCell(ratingString(f.Rating)).SetTextColor(tview.Styles.PrimaryTextColor)
			case 1:
				cell = tview.NewTableCell(f.Title).SetTextColor(tview.Styles.PrimaryTextColor)
			case 2:
				cell = tview.NewTableCell(f.Author).SetTextColor(tview.Styles.PrimaryTextColor)
			case 3:
				cell = tview.NewTableCell(f.Date).SetTextColor(tview.Styles.PrimaryTextColor)
			default:
				cell = tview.NewTableCell("").SetTextColor(tview.Styles.PrimaryTextColor)
			}

			browser.list.SetCell(r, c, cell)
		}
	}
}

func (browser *IdgamesBrowser) populateDetails(idgame goidgames.Idgame) {
	browser.fileDetailsText.Clear()
	fmt.Fprintf(browser.fileDetailsText, "%s", idgame.Filename)
}

func ratingString(rating float32) string {
	return strings.Repeat("*", int(rating)) + strings.Repeat("-", 5-int(rating))
}
