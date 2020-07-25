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

	fixRows, fixCols := 1, 4
	rows, cols := len(allGames), 0
	if config.LegacyModList {
		cols = games.MaxModCount() - 1
	}

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
					cell = tview.NewTableCell(game.RatingString()).SetTextColor(tview.Styles.TitleColor)
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

			// get out
			case 'q':
				app.Stop()

			// show credits and license
			case 'c':
				// c again to toggle
				frontPage, _ := detailSidePagesSub2.GetFrontPage()
				if frontPage == pageLicense {
					appModeNormal()
					return nil
				}
				lp := makeLicense()
				detailSidePagesSub2.AddPage(pageLicense, lp, true, true)
				detailSidePagesSub2.SwitchToPage(pageLicense)
				app.SetFocus(lp)
				return nil

			// options
			case 'o':
				optionsDiag := makeOptions()
				contentPages.AddPage(pageOptions, optionsDiag, true, false)
				contentPages.SwitchToPage(pageOptions)
				app.SetFocus(optionsDiag)

			// open dialog to insert new game
			case 'i':
				newForm := makeAddEditGame(nil)
				detailPages.AddPage(pageAddEdit, newForm, true, false)
				detailPages.SwitchToPage(pageAddEdit)
				app.SetFocus(newForm)
				return nil

			// increase game rating
			case '+':
				allGames[r-fixRows].Rate(1)
				c := tview.NewTableCell(allGames[r-fixRows].RatingString()).SetTextColor(tview.Styles.SecondaryTextColor)
				gamesTable.SetCell(r, 0, c)
				games.Persist()

			// decrease game rating
			case '-':
				allGames[r-fixRows].Rate(-1)
				c := tview.NewTableCell(allGames[r-fixRows].RatingString()).SetTextColor(tview.Styles.SecondaryTextColor)
				gamesTable.SetCell(r, 0, c)
				games.Persist()

			// open dialog to add mod to game
			case 'a':
				if r > 0 {
					mtm := makeModTreeMaker(&allGames[r-fixRows])
					modTree := mtm()
					detailSidePagesSub2.AddPage(pageModSelector, modTree, true, false)
					detailSidePagesSub2.SwitchToPage(pageModSelector)
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
								selectedGameChanged(&allGames[r-fixRows])
								games.Persist()
							}
						}
					}

					if config.DeleteWithoutWarning {
						removeMod()
						return nil
					}
					g := allGames[r-fixRows]
					contentPages.AddPage(pageYouSure, makeYouSureBox(fmt.Sprintf(deleteModQuestion, g.Mods[len(g.Mods)-1], g.Name), removeMod, 2, r+2), true, true)
					return nil

				}
				return nil

			case 'e':
				if r > 0 {
					addEdit := makeAddEditGame(&allGames[r-fixRows])
					detailPages.AddPage(pageAddEdit, addEdit, true, false)
					detailPages.SwitchToPage(pageAddEdit)
					app.SetFocus(addEdit)
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

			if config.DeleteWithoutWarning {
				remove()
				return nil
			}

			g := allGames[r-fixRows]
			contentPages.AddPage(pageYouSure, makeYouSureBox(fmt.Sprintf(deleteGameQuestion, g.Name), remove, 2, r+2), true, true)
			return nil
		}

		return event
	})
}
