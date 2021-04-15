package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/goidgames"
)

//"github.com/zmnpl/goidgames"

type IdgamesBrowser struct {
	app         *tview.Application
	layout      *tview.Grid
	list        *tview.Table
	fileDetails *tview.TextView
	reviews     *tview.TextView
	search      *tview.InputField
	idgames     []goidgames.Idgame
}

func NewIdgamesBrowser(app *tview.Application) *IdgamesBrowser {
	browser := &IdgamesBrowser{app: app}

	layout := tview.NewGrid()
	layout.SetRows(5, -1, 5)
	layout.SetColumns(-1, -1)

	browser.layout = layout

	browser.initList()
	browser.initDetails()
	browser.initSearchForm()

	return browser
}

func (b *IdgamesBrowser) initSearchForm() {
	searchForm := tview.NewForm()
	searchForm.SetHorizontal(true).SetBorder(true)

	search := tview.NewInputField().SetLabel("Search Query").SetText("")
	searchForm.AddFormItem(search)

	searchByTitle := true
	searchByAuthor := false
	searchForm.AddCheckbox("By Title", true, func(checked bool) { searchByTitle = checked })
	searchForm.AddCheckbox("By Author", false, func(checked bool) { searchByAuthor = checked })

	searchForm.AddButton("Search", func() {
		query := search.GetText()
		if len(query) == 0 {
			b.UpdateLatest()
		} else {
			types := make([]string, 0)
			if searchByTitle {
				types = append(types, goidgames.SEARCH_TYPE_TITLE)
			}
			if searchByAuthor {
				types = append(types, goidgames.SEARCH_TYPE_AUTHOR)
			}

			b.UpdateSearch(search.GetText(), types)
		}
		app.SetFocus(b.list)
	})

	b.layout.AddItem(searchForm, 0, 0, 1, 2, 0, 0, true)

	b.search = search
}

func (b *IdgamesBrowser) initDetails() {
	details := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	details.SetBorder(true).
		SetBorderPadding(0, 0, 1, 1)

	b.layout.AddItem(details, 1, 1, 1, 1, 0, 0, false)

	details.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyTAB {
			b.app.SetFocus(b.search)
			return nil
		}
		if k == tcell.KeyBacktab {
			b.app.SetFocus(b.list)
			return nil
		}
		return event
	})

	b.fileDetails = details
}

func (b *IdgamesBrowser) initList() {
	list := tview.NewTable().
		SetFixed(1, 2).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')
	list.SetBorder(true)

	b.layout.AddItem(list, 1, 0, 1, 1, 0, 0, false)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyTAB {
			b.app.SetFocus(b.fileDetails)
			return nil
		}
		if k == tcell.KeyBacktab {
			b.app.SetFocus(b.search)
			return nil
		}
		return event
	})

	list.SetSelectedFunc(func(r int, c int) {
		switch {
		case r > 0:
			// TODO - what to do now? select multiple or install right away?
			//fmt.Println(b.idgames[r-1].Url)
		}
	})

	b.list = list
}

// updateGameDetails iterates the given slice and fetches the detail data from Idgames via the api's get function
func updateGameDetails(idgames []goidgames.Idgame) {
	for i := range idgames {
		g, err := goidgames.Get(idgames[i].Id, "")
		if err != nil {
			continue
		}
		idgames[i] = g
	}
}

func (browser *IdgamesBrowser) UpdateSearch(query string, types []string) {
	go func() {
		app.QueueUpdateDraw(func() {
			idgames, _ := goidgames.SearchMultipleTypes(query, types, goidgames.SEARCH_SORT_RATING, goidgames.SEARCH_SORT_DESC)

			go func() {
				updateGameDetails(idgames)
			}()

			browser.populateList(idgames)
		})
	}()
}

func (browser *IdgamesBrowser) UpdateLatest() {
	go func() {
		app.QueueUpdateDraw(func() {
			idgames, _ := goidgames.LatestFiles(10, 0)

			go func() {
				updateGameDetails(idgames)
			}()

			browser.populateList(idgames)
		})
	}()
}

func (browser *IdgamesBrowser) populateList(idgames []goidgames.Idgame) {
	browser.list.Clear()
	browser.idgames = idgames

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
	browser.list.ScrollToBeginning()
}

func (browser *IdgamesBrowser) populateDetails(idgame goidgames.Idgame) {
	browser.fileDetails.Clear()
	fmt.Fprintf(browser.fileDetails, "%s", idgame.Textfile)
}

func ratingString(rating float32) string {
	return strings.Repeat("*", int(rating)) + strings.Repeat("-", 5-int(rating))
}
