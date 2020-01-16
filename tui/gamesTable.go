package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	gameTableHeaderName       = "Name"
	gameTableHeaderSourcePort = "SourcePort"
	gameTableHeaderIwad       = "Iwad"
	gameTableHeaderMods       = "Mods"
)

// center table with mods
func makeGamesTable() *tview.Table {
	gamesTable = tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')

	return gamesTable
}

func populateGamesTable() {
	gamesTable.Clear()
	allGames := games.GetInstance()

	rows, cols := len(allGames), games.MaxModCount()-1
	fixRows, fixCols := 1, 3

	for r := 0; r < rows+fixRows; r++ {
		var game games.Game
		if r > 0 {
			game = allGames[r-fixRows]
		}
		for c := 0; c < cols+fixCols; c++ {
			var cell *tview.TableCell

			if r < 1 {
				switch c {
				case 0:
					cell = tview.NewTableCell(gameTableHeaderName).SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell(gameTableHeaderSourcePort).SetTextColor(tview.Styles.SecondaryTextColor)
				case 2:
					cell = tview.NewTableCell(gameTableHeaderIwad).SetTextColor(tview.Styles.SecondaryTextColor)
				case 3:
					cell = tview.NewTableCell(gameTableHeaderMods).SetTextColor(tview.Styles.SecondaryTextColor)
				default:
					cell = tview.NewTableCell("").SetTextColor(tview.Styles.SecondaryTextColor)
				}
			} else {
				switch c {
				case 0:
					cell = tview.NewTableCell(game.Name).SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell(game.SourcePort).SetTextColor(tview.Styles.PrimaryTextColor)
				case 2:
					cell = tview.NewTableCell(game.Iwad).SetTextColor(tview.Styles.PrimaryTextColor)
				default:
					i := c - fixCols
					if i < len(game.Mods) {
						cell = tview.NewTableCell(game.Mods[i]).SetTextColor(tview.Styles.PrimaryTextColor)
					} else {
						cell = tview.NewTableCell("").SetTextColor(tview.Styles.PrimaryTextColor)
					}
				}
			}
			gamesTable.SetCell(r, c, cell)
		}
	}

	makeModTreeMaker := func(selectedGame *games.Game) func() *tview.TreeView {
		return func() *tview.TreeView {
			return makeModTree(selectedGame)
		}
	}
	modTreeMaker := makeModTreeMaker(&games.Game{})

	//makeCellPulser := func(c *tview.TableCell) func() {
	//	return func() {
	//		r, g, b := tcell.ColorOrange.RGB()
	//		c.BackgroundColor = tcell.NewRGBColor(r, g, b)
	//		for i := 0; i <= 1000; i += 16 {
	//			time.Sleep(16 * time.Millisecond)
	//			r = r * 2
	//			g = g * 2
	//			b = b * 2
	//			c.BackgroundColor = tcell.NewRGBColor(r, g, b)
	//			app.Draw()
	//		}
	//	}
	//}
	//cellPulser := makeCellPulser(tview.NewTableCell(""))

	gamesTable.SetSelectionChangedFunc(func(r int, c int) {
		var g *games.Game
		//var cell *tview.TableCell
		switch r {
		case 0:
			g = &games.Game{}
			//cell = tview.NewTableCell("")
		default:
			g = &allGames[r-fixRows]
			//cell = gamesTable.GetCell(r, len(g.Mods)+3)
		}
		selectedGameChanged(g)
		modTreeMaker = makeModTreeMaker(g)
		//cellPulser = makeCellPulser(cell)
	})

	gamesTable.SetSelectedFunc(func(r int, c int) {
		switch {
		case r > 0:
			allGames[r-fixRows].Run()
		}
	})

	gamesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		r, _ := gamesTable.GetSelection()

		if k == tcell.KeyRune {
			switch event.Rune() {
			// open dialog to add mod to game
			case 'a':
				modTree := modTreeMaker()
				actionPager.AddPage(pageModSelector, modTree, true, false)
				actionPager.SwitchToPage(pageModSelector)
				app.SetFocus(modTree)
				return nil

			// remove last mod from game
			case 'r':
				if r > 0 {
					mods := allGames[r-fixRows].Mods
					if len(mods) > 0 {
						allGames[r-fixRows].Mods = mods[:len(mods)-1]
						populateGamesTable()
						games.Persist()
					}
				}
				return nil

			// open dialog to insert new game
			case 'i':
				actionPager.SwitchToPage(pageNewForm)
				app.SetFocus(newForm)
				return nil

			// show credits and license
			case 'c':
				frontPage, _ := actionPager.GetFrontPage()
				if frontPage == pageLicense {
					appModeNormal()
					return nil
				}
				actionPager.SwitchToPage(pageLicense)
				app.SetFocus(licensePage)
				return nil
			}
		}

		if k == tcell.KeyDelete && r > 0 {
			if r == gamesTable.GetRowCount()-1 {
				gamesTable.Select(r-1, 0)
			}
			games.RemoveGameAt(r - fixRows)
			return nil
		}

		return event
	})
}
