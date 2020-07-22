package tui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	gameTableHeaderRating     = "Rating"
	gameTableHeaderName       = "Name"
	gameTableHeaderSourcePort = "SourcePort"
	gameTableHeaderIwad       = "Iwad"
	gameTableHeaderMods       = "Mods"

	deleteGameQuestion = "Delete '%v'?"
	deleteModQuestion  = "Remove '%v' from '%v'?"
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
	fixRows, fixCols := 1, 4

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
					cell = tview.NewTableCell(gameTableHeaderRating).SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell(gameTableHeaderName).SetTextColor(tview.Styles.SecondaryTextColor)
				case 2:
					cell = tview.NewTableCell(gameTableHeaderSourcePort).SetTextColor(tview.Styles.SecondaryTextColor)
				case 3:
					cell = tview.NewTableCell(gameTableHeaderIwad).SetTextColor(tview.Styles.SecondaryTextColor)
				case 4:
					cell = tview.NewTableCell(gameTableHeaderMods).SetTextColor(tview.Styles.SecondaryTextColor)
				default:
					cell = tview.NewTableCell("").SetTextColor(tview.Styles.SecondaryTextColor)
				}
			} else {
				switch c {
				case 0:
					cell = tview.NewTableCell(game.RatingString()).SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell(game.Name).SetTextColor(tview.Styles.SecondaryTextColor)
				case 2:
					cell = tview.NewTableCell(game.SourcePort).SetTextColor(tview.Styles.PrimaryTextColor)
				case 3:
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
			case '+':
				allGames[r-fixRows].Rate(1)
				c := tview.NewTableCell(allGames[r-fixRows].RatingString()).SetTextColor(tview.Styles.SecondaryTextColor)
				gamesTable.SetCell(r, 0, c)
				games.Persist()
			case '-':
				allGames[r-fixRows].Rate(-1)
				c := tview.NewTableCell(allGames[r-fixRows].RatingString()).SetTextColor(tview.Styles.SecondaryTextColor)
				gamesTable.SetCell(r, 0, c)
				games.Persist()

			case 'o':
				optionsDiag := makeOptions()
				bigMainPager.AddPage(pageOptions, optionsDiag, true, false)
				bigMainPager.SwitchToPage(pageOptions)
				app.SetFocus(optionsDiag)
			// open dialog to add mod to game
			case 'a':
				if r > 0 {
					mtm := makeModTreeMaker(&allGames[r-fixRows])
					modTree := mtm()
					actionPager.AddPage(pageModSelector, modTree, true, false)
					actionPager.SwitchToPage(pageModSelector)
					app.SetFocus(modTree)
					return nil
				}

			// remove last mod from game
			case 'r':
				mods := allGames[r-fixRows].Mods
				if len(mods) > 0 {
					removeMod := func() {
						if r > 0 {
							if len(mods) > 0 {
								allGames[r-fixRows].Mods = mods[:len(mods)-1]
								populateGamesTable()
								games.Persist()
							}
						}
					}

					if config.WarnBeforeDelete {
						g := allGames[r-fixRows]
						bigMainPager.AddPage(pageYouSure, makeYouSureBox(fmt.Sprintf(deleteModQuestion, g.Mods[len(g.Mods)-1], g.Name), removeMod, 2, r+2), true, true)
						return nil
					}

					removeMod()
				}
				return nil

			// open dialog to insert new game
			case 'i':
				newForm := makeAddEditGame(nil)
				actionPager.AddPage(pageNewForm, newForm, true, false)
				actionPager.SwitchToPage(pageNewForm)
				app.SetFocus(newForm)
				return nil

			case 'e':
				if r > 0 {
					customParameters := makeAddEditGame(&allGames[r-fixRows])
					actionPager.AddPage(pageParamsEdit, customParameters, true, false)
					actionPager.SwitchToPage(pageParamsEdit)
					app.SetFocus(customParameters)
					return nil
				}

			case 'd':
				if r > 0 {
					gameOverview := makeModList(&allGames[r-fixRows])
					actionPager.AddPage(pageGameOverview, gameOverview, true, false)
					actionPager.SwitchToPage(pageGameOverview)
					return nil
				}

			case 's':
				games.SortAlph()
				populateGamesTable()
				return nil
			}
		}

		if k == tcell.KeyDelete && r > 0 {
			remove := func() {
				if r == gamesTable.GetRowCount()-1 {
					gamesTable.Select(r-1, 0)
				}
				games.RemoveGameAt(r - fixRows)
			}

			if config.WarnBeforeDelete {
				g := allGames[r-fixRows]
				bigMainPager.AddPage(pageYouSure, makeYouSureBox(fmt.Sprintf(deleteGameQuestion, g.Name), remove, 2, r+2), true, true)
				return nil
			}

			remove()

			return nil
		}

		return event
	})
}
